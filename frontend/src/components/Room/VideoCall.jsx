import { useState } from 'react';
import {
  SpeakerWaveIcon,
  SpeakerXMarkIcon,
  VideoCameraIcon,
  VideoCameraSlashIcon,
  UserGroupIcon,
  PhoneXMarkIcon
} from '@heroicons/react/24/solid';

const VideoMeetingRoom = () => {
  const [isMicOn, setIsMicOn] = useState(true);
  const [isVideoOn, setIsVideoOn] = useState(true);
  const meetingName = "Project Sync"; // This would come from props or context

  const toggleMic = () => setIsMicOn(!isMicOn);
  const toggleVideo = () => setIsVideoOn(!isVideoOn);
  const handleLeave = () => {
    // Handle leave meeting logic here
    console.log('Leaving meeting...');
  };

  return (
    <div className="h-screen bg-gray-900 flex flex-col">
      {/* Header */}
      <header className="bg-gray-800 p-4 flex items-center justify-between">
        <h1 className="text-white text-xl font-semibold">{meetingName}</h1>
        <div className="flex items-center space-x-6">
          <div className="flex items-center text-gray-300">
            <UserGroupIcon className="w-5 h-5 mr-2" />
            <span>4</span>
          </div>
          <button
            onClick={handleLeave}
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
            {isVideoOn ? (
              <div className="absolute bottom-4 left-4 bg-gray-900 px-2 py-1 rounded text-white text-sm">
                You
              </div>
            ) : (
              <div className="h-full flex items-center justify-center">
                <div className="w-20 h-20 rounded-full bg-gray-600 flex items-center justify-center text-white text-2xl">
                  Y
                </div>
              </div>
            )}
          </div>
          
          {/* Other participants */}
          <div className="bg-gray-800 rounded-lg relative">
            <div className="h-full flex items-center justify-center">
              <div className="w-20 h-20 rounded-full bg-gray-600 flex items-center justify-center text-white text-2xl">
                J
              </div>
            </div>
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
              isMicOn ? 'bg-gray-600 hover:bg-gray-700' : 'bg-red-600 hover:bg-red-700'
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
              isVideoOn ? 'bg-gray-600 hover:bg-gray-700' : 'bg-red-600 hover:bg-red-700'
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