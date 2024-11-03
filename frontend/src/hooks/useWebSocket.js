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
  var [peerConnection, setPeerConnection] = useState(null);
  const socket = useRef(null);

  useEffect(() => {
    startCall()
    // Initialize WebSocket connection
    const signalingServerUrl = `${import.meta.env.VITE_WS_URI}/api/room/v1/${roomID}/ws`;
    socket.current = new WebSocket(signalingServerUrl);
    socket.current.onopen = () => {
      console.log("Connected to signaling server");
    };
    socket.current.onmessage = (message) => {
      const data = JSON.parse(message.data);
      console.log(data)
      handleSignalingData(data);
    };

    return () => {
      socket.current.close();
    };
  }, []);

  const startCall = async () => {
    const pc = new RTCPeerConnection(configuration);
    setPeerConnection(pc);

    const stream = await navigator.mediaDevices.getUserMedia({
      video: isVideoOn,
      audio: isMicOn,
    });
  
    setLocalStream(stream);
    if (localVideoRef.current) {
      localVideoRef.current.srcObject = stream;
    }
  
    stream.getTracks().forEach((track) => pc.addTrack(track, stream));
  
    pc.onicecandidate = (event) => {
      if (event.candidate) {
        socket.current.send(JSON.stringify({ type: "ice-candidate", candidate: event.candidate }));
      }
    };
  
    pc.ontrack = (event) => {
      if (remoteVideoRef.current) {
        remoteVideoRef.current.srcObject = event.streams[0];
      }
    };
  
    const offer = await pc.createOffer();
    await pc.setLocalDescription(offer);
    socket.current.send(JSON.stringify({ type: "offer", offer }));
  };
  

  const handleSignalingData = (data) => {
    switch (data.type) {
      case "offer": {
        const offer = new RTCSessionDescription(data.offer);
        peerConnection.setRemoteDescription(offer)
          .then(() => {
            return peerConnection.createAnswer();
          })
          .then((answer) => {
            return peerConnection.setLocalDescription(answer);
          })
          .then(() => {
            socket.current.send(JSON.stringify({ type: "answer", answer: peerConnection.localDescription }));
          })
          .catch((error) => {
            console.error("Error handling offer:", error);
          });
        break;
      }
      case "answer": {
        const answer = new RTCSessionDescription(data.answer);
        if (peerConnection) {
          peerConnection.setRemoteDescription(answer).catch((error) => {
            console.error("Error setting remote description:", error);
          });
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
