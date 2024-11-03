import { useEffect, useRef, useState } from "react";

const configuration = {
  iceServers: [
    { urls: "stun:stun.l.google.com:19302" }, // STUN server
  ],
};

const useWebRTC = (isVideoOn, isMicOn, roomID) => {
  const localVideoRef = useRef(null);
  const remoteVideoRef = useRef(null);
  const [localStream, setLocalStream] = useState(null);
  const [peerConnection, setPeerConnection] = useState(null);
  const socket = useRef(null);

  useEffect(() => {
    // Initialize WebSocket connection
    const signalingServerUrl = `${
      import.meta.env.VITE_WS_URI
    }/api/room/v1/${roomID}/ws`;
    socket.current = new WebSocket(signalingServerUrl);

    socket.current.onopen = () => {
      console.log("Connected to signaling server");
      // Only start the call after the WebSocket is fully connected
      startCall();
    };

    socket.current.onmessage = (message) => {
      const data = JSON.parse(message.data);
      console.log(data);
      handleSignalingData(data);
    };

    return () => {
      socket.current.close();
    };
  }, [roomID]); // eslint-disable-line react-hooks/exhaustive-deps

  const startCall = async () => {
    try {
      // Create the peer connection locally in the function
      const pc = new RTCPeerConnection(configuration);

      // Update state with the peer connection
      setPeerConnection(pc);

      const stream = await navigator.mediaDevices.getUserMedia({
        video: isVideoOn,
        audio: isMicOn,
      });

      setLocalStream(stream);
      if (localVideoRef.current) {
        localVideoRef.current.srcObject = stream;
      }

      // Add local tracks to the peer connection
      stream.getTracks().forEach((track) => pc.addTrack(track, stream));

      // Set up ICE candidate handling
      pc.onicecandidate = (event) => {
        if (event.candidate) {
          socket.current.send(
            JSON.stringify({
              type: "ice-candidate",
              candidate: event.candidate,
            })
          );
        }
      };

      // Set up handling of incoming tracks
      pc.ontrack = (event) => {
        if (remoteVideoRef.current) {
          remoteVideoRef.current.srcObject = event.streams[0];
        }
      };

      // Create offer, set local description, and send it via signaling
      const offer = await pc.createOffer();
      await pc.setLocalDescription(offer);
      socket.current.send(JSON.stringify({ type: "offer", offer }));
    } catch (error) {
      console.log("Error while initiating call", error);
    }
  };

  const handleSignalingData = (data) => {
    switch (data.type) {
      case "offer": {
        const offer = new RTCSessionDescription(data.offer);
        if (peerConnection.signalingState === "stable") {
          // Only set remote description if stable
          peerConnection
            .setRemoteDescription(offer)
            .then(() => {
              return peerConnection.createAnswer();
            })
            .then((answer) => {
              return peerConnection.setLocalDescription(answer);
            })
            .then(() => {
              socket.current.send(
                JSON.stringify({
                  type: "answer",
                  answer: peerConnection.localDescription,
                })
              );
            })
            .catch((error) => {
              console.error("Error handling offer:", error);
            });
        } else {
          console.warn(
            "Received an offer but the connection is not stable, ignoring it."
          );
        }
        break;
      }
      case "answer": {
        const answer = new RTCSessionDescription(data.answer);
        if (peerConnection.signalingState === "have-local-offer") {
          // Check state before setting answer
          peerConnection.setRemoteDescription(answer).catch((error) => {
            console.error("Error setting remote description:", error);
          });
        } else {
          console.warn(
            "Received an answer but the connection is not in a state to accept it, ignoring it."
          );
        }
        break;
      }
      case "ice-candidate": {
        const candidate = new RTCIceCandidate(data.candidate);
        if (peerConnection) {
          peerConnection.addIceCandidate(candidate).catch((error) => {
            console.error("Error adding received ice candidate:", error);
          });
        }
        break;
      }
      default:
        break;
    }
  };

  const leaveCall = () => {
    if (peerConnection) {
      peerConnection.close();
      setPeerConnection(null);
      if (localStream) {
        localStream.getTracks().forEach((track) => track.stop());
        setLocalStream(null);
      }
    }
  };

  return {
    localVideoRef,
    remoteVideoRef,
    leaveCall,
  };
};

export default useWebRTC;
