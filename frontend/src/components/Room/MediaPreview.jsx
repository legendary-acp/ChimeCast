import { useState, useEffect, useRef } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { 
  SpeakerWaveIcon, 
  SpeakerXMarkIcon,
  VideoCameraIcon, 
  VideoCameraSlashIcon,
} from '@heroicons/react/24/solid';


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
  const navigate = useNavigate();

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
  const JoinRoom = () => {
    console.log(`Joining room with ID: ${room.id}`);
    navigate(`/meeting/${room.id}`, { state: { room: room } });
  };

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gray-800 p-4">
      <div className="bg-gray-900 p-6 rounded-lg shadow-lg text-white space-y-4">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-xl font-bold">
              Media Preview for: {room.name}
            </h2>
          </div>
          <button
            onClick={JoinRoom}
            className="flex px-4 py-2 rounded-full text-white bg-blue-700 hover:bg-green-500 hover:scale-105"
          >
            Join
          </button>
        </div>
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
              <div className="w-20 h-20 rounded-full bg-orange-600 flex items-center justify-center text-white text-2xl">
                J
              </div>
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
            className={`p-4 rounded-full ${
              hasVideo ? 'bg-gray-600 hover:bg-gray-700' : 'bg-red-600 hover:bg-red-700'
            }`}
          >
            {hasVideo ? (
              <VideoCameraIcon className="w-6 h-6 text-white" />
            ) : (
              <VideoCameraSlashIcon className="w-6 h-6 text-white" />
            )}
          </button>

          <button
            onClick={toggleAudio}
            className={`p-4 rounded-full ${
              hasAudio ? 'bg-gray-600 hover:bg-gray-700' : 'bg-red-600 hover:bg-red-700'
            }`}
          >
            {hasAudio ? (
              <SpeakerWaveIcon className="w-6 h-6 text-white" />
            ) : (
              <SpeakerXMarkIcon className="w-6 h-6 text-white" />
            )}
          </button>
        </div>
      </div>
    </div>
  );
};

export default MediaPreview;
