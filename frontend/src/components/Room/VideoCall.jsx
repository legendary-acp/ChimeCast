import { useState } from "react";
import { useLocation } from "react-router-dom";
import {
  SpeakerWaveIcon,
  SpeakerXMarkIcon,
  VideoCameraIcon,
  VideoCameraSlashIcon,
  UserGroupIcon,
  PhoneXMarkIcon,
} from "@heroicons/react/24/solid";
import useWebRTC from "../../hooks/useWebSocket";

const VideoMeetingRoom = () => {
  const [isMicOn, setIsMicOn] = useState(true);
  const [isVideoOn, setIsVideoOn] = useState(true);
  const room = useLocation().state?.room;
  const { localVideoRef, remoteVideoRef, leaveCall } = useWebRTC(
    isVideoOn,
    isMicOn,
    room.id
  );

  const toggleVideo = () => {
    if (localVideoRef.current && localVideoRef.current.srcObject) {
      localVideoRef.current.srcObject
        .getVideoTracks()
        .forEach((track) => (track.enabled = !track.enabled));
      setIsVideoOn(!isVideoOn);
    }
  };

  const toggleMic = () => {
    if (localVideoRef.current && localVideoRef.current.srcObject) {
      localVideoRef.current.srcObject
        .getAudioTracks()
        .forEach((track) => (track.enabled = !track.enabled));
      setIsMicOn(!isMicOn);
    }
  };

  return (
    <div className="h-screen bg-gray-900 flex flex-col">
      {/* Header */}
      <header className="bg-gray-800 p-4 flex items-center justify-between">
        <h1 className="text-white text-xl font-semibold">{room.name}</h1>
        <div className="flex items-center space-x-6">
          <div className="flex items-center text-gray-300">
            <UserGroupIcon className="w-5 h-5 mr-2" />
            <span>2</span>
          </div>
          <button
            onClick={leaveCall}
            className="px-3 py-2 bg-red-600 hover:bg-red-700 rounded-lg flex items-center justify-center transition-colors"
          >
            <PhoneXMarkIcon className="w-5 h-5 text-white" />
            <span className="ml-2 text-white text-sm font-medium">Leave</span>
          </button>
        </div>
      </header>

      {/* Main content */}
      <main className="flex-1 p-4 relative">
        {/* Video grid */}
        <div className="grid grid-cols-2 gap-4 h-full">
          {/* Main video */}
          <div className="bg-gray-800 rounded-lg relative">
            <video
              ref={localVideoRef}
              autoPlay
              muted
              className="w-full h-full rounded-lg"
            />
            {!isVideoOn && (
              <div className="absolute inset-0 flex items-center justify-center text-gray-400">
                <div className="w-20 h-20 rounded-full bg-orange-600 flex items-center justify-center text-white text-2xl">
                  P
                </div>
              </div>
            )}
            <div className="absolute bottom-4 left-4 bg-gray-900 px-2 py-1 rounded text-white text-sm">
              You
            </div>
          </div>

          {/* Other participants */}
          <div className="bg-gray-800 rounded-lg relative">
            <video
              ref={remoteVideoRef}
              autoPlay
              className="w-full h-full rounded-lg"
            />
            <div className="absolute bottom-4 left-4 bg-gray-900 px-2 py-1 rounded text-white text-sm">
              John
            </div>
          </div>
        </div>
      </main>

      {/* Controls */}
      <footer className="bg-gray-800 p-4">
        <div className="flex justify-center items-center space-x-4">
          {/* Audio control */}
          <button
            onClick={toggleMic}
            className={`p-4 rounded-full ${
              isMicOn
                ? "bg-gray-600 hover:bg-gray-700"
                : "bg-red-600 hover:bg-red-700"
            } transition-colors`}
          >
            {isMicOn ? (
              <SpeakerWaveIcon className="w-6 h-6 text-white" />
            ) : (
              <SpeakerXMarkIcon className="w-6 h-6 text-white" />
            )}
          </button>

          {/* Video control */}
          <button
            onClick={toggleVideo}
            className={`p-4 rounded-full ${
              isVideoOn
                ? "bg-gray-600 hover:bg-gray-700"
                : "bg-red-600 hover:bg-red-700"
            } transition-colors`}
          >
            {isVideoOn ? (
              <VideoCameraIcon className="w-6 h-6 text-white" />
            ) : (
              <VideoCameraSlashIcon className="w-6 h-6 text-white" />
            )}
          </button>
        </div>
      </footer>
    </div>
  );
};

export default VideoMeetingRoom;
