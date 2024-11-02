import { useState, useEffect, useRef } from "react";
import { useLocation } from "react-router-dom";

const MediaPreview = () => {
  const [hasVideo, setHasVideo] = useState(true);
  const [hasAudio, setHasAudio] = useState(true);
  const [audioLevel, setAudioLevel] = useState(0);
  const [error, setError] = useState("");
  const videoRef = useRef(null);
  const audioContextRef = useRef(null);
  const analyserRef = useRef(null);
  const animationFrameRef = useRef(null);

  const room = useLocation().state?.room;

  useEffect(() => {
    startMedia();

    return () => {
      if (animationFrameRef.current) {
        cancelAnimationFrame(animationFrameRef.current);
      }
      if (audioContextRef.current) {
        audioContextRef.current.close();
      }
    };
  }, []);

  const startMedia = async () => {
    try {
      // Requesting both audio and video permissions
      const stream = await navigator.mediaDevices.getUserMedia({
        video: true,
        audio: true,
      });

      // Set video source
      if (videoRef.current) {
        videoRef.current.srcObject = stream;
      }

      // Setup audio context and analyser for audio level visualization
      audioContextRef.current = new AudioContext();
      const source = audioContextRef.current.createMediaStreamSource(stream);
      analyserRef.current = audioContextRef.current.createAnalyser();
      source.connect(analyserRef.current);

      // Start analyzing audio levels
      const updateAudioLevel = () => {
        const dataArray = new Uint8Array(analyserRef.current.frequencyBinCount);
        analyserRef.current.getByteFrequencyData(dataArray);
        const average = dataArray.reduce((a, b) => a + b) / dataArray.length;
        setAudioLevel(average);
        animationFrameRef.current = requestAnimationFrame(updateAudioLevel);
      };
      updateAudioLevel();

      // Clear any previous errors
      setError("");
    } catch (err) {
      // Handle errors if permissions are denied
      setError(
        "Failed to access audio and video devices. Please grant permissions."
      );
      console.error("Failed to access media devices:", err);
    }
  };

  const toggleVideo = () => {
    if (videoRef.current && videoRef.current.srcObject) {
      videoRef.current.srcObject
        .getVideoTracks()
        .forEach((track) => (track.enabled = !track.enabled));
      setHasVideo((prev) => !prev);
    }
  };

  const toggleAudio = () => {
    if (videoRef.current && videoRef.current.srcObject) {
      videoRef.current.srcObject
        .getAudioTracks()
        .forEach((track) => (track.enabled = !track.enabled));
      setHasAudio((prev) => !prev);
    }
  };

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gray-800 p-4">
      <div className="bg-gray-900 p-6 rounded-lg shadow-lg text-white space-y-4">
        <h2 className="text-xl font-bold">Media Preview</h2>
        <h3 className="text-xl font-bold">Joining {room.name}</h3>

        {/* Display any error messages */}
        {error && <div className="text-red-500">{error}</div>}

        {/* Video Preview */}
        <div className="relative aspect-video bg-gray-700 rounded-lg overflow-hidden">
          <video
            ref={videoRef}
            autoPlay
            playsInline
            muted
            className="w-full h-full object-cover rounded-lg"
          />
          {!hasVideo && (
            <div className="absolute inset-0 flex items-center justify-center text-gray-400">
              <span>Video is Off</span>
            </div>
          )}
        </div>

        {/* Audio Level Indicator */}
        <div className="text-center">
          Audio Level:
          <div className="w-full bg-gray-700 rounded-full h-2 mt-1">
            <div
              className="bg-green-500 h-2 rounded-full"
              style={{ width: `${Math.min(audioLevel, 100)}%` }}
            ></div>
          </div>
        </div>

        {/* Controls */}
        <div className="flex space-x-4 justify-center mt-4">
          <button
            onClick={toggleVideo}
            className={`px-4 py-2 rounded-full text-white ${
              hasVideo ? "bg-red-500" : "bg-green-500"
            }`}
          >
            {hasVideo ? (
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth="1.5"
                stroke="currentColor"
                className="size-6"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="m15.75 10.5 4.72-4.72a.75.75 0 0 1 1.28.53v11.38a.75.75 0 0 1-1.28.53l-4.72-4.72M12 18.75H4.5a2.25 2.25 0 0 1-2.25-2.25V9m12.841 9.091L16.5 19.5m-1.409-1.409c.407-.407.659-.97.659-1.591v-9a2.25 2.25 0 0 0-2.25-2.25h-9c-.621 0-1.184.252-1.591.659m12.182 12.182L2.909 5.909M1.5 4.5l1.409 1.409"
                />
              </svg>
            ) : (
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth="1.5"
                stroke="currentColor"
                className="size-6"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="m15.75 10.5 4.72-4.72a.75.75 0 0 1 1.28.53v11.38a.75.75 0 0 1-1.28.53l-4.72-4.72M4.5 18.75h9a2.25 2.25 0 0 0 2.25-2.25v-9a2.25 2.25 0 0 0-2.25-2.25h-9A2.25 2.25 0 0 0 2.25 7.5v9a2.25 2.25 0 0 0 2.25 2.25Z"
                />
              </svg>
            )}
          </button>

          <button
            onClick={toggleAudio}
            className={`px-4 py-2 rounded-full text-white ${
              hasAudio ? "bg-red-500" : "bg-green-500"
            }`}
          >
            {hasAudio ? (
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth="1.5"
                stroke="currentColor"
                className="size-6"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M17.25 9.75 19.5 12m0 0 2.25 2.25M19.5 12l2.25-2.25M19.5 12l-2.25 2.25m-10.5-6 4.72-4.72a.75.75 0 0 1 1.28.53v15.88a.75.75 0 0 1-1.28.53l-4.72-4.72H4.51c-.88 0-1.704-.507-1.938-1.354A9.009 9.009 0 0 1 2.25 12c0-.83.112-1.633.322-2.396C2.806 8.756 3.63 8.25 4.51 8.25H6.75Z"
                />
              </svg>
            ) : (
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth="1.5"
                stroke="currentColor"
                className="size-6"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M19.114 5.636a9 9 0 0 1 0 12.728M16.463 8.288a5.25 5.25 0 0 1 0 7.424M6.75 8.25l4.72-4.72a.75.75 0 0 1 1.28.53v15.88a.75.75 0 0 1-1.28.53l-4.72-4.72H4.51c-.88 0-1.704-.507-1.938-1.354A9.009 9.009 0 0 1 2.25 12c0-.83.112-1.633.322-2.396C2.806 8.756 3.63 8.25 4.51 8.25H6.75Z"
                />
              </svg>
            )}
          </button>
        </div>
      </div>
    </div>
  );
};

export default MediaPreview;
