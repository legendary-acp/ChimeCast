import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import PropTypes from "prop-types";
import { Disclosure, Menu } from "@headlessui/react";
import {
  fetchRooms,
  getRoomStatus,
  joinRoom,
} from "../../services/roomService";
import { logout } from "../../services/authService";
import CreateRoomPanel from "../Room/CreateRoom";
import formatTimestamp from "../../utils/dateUtils";
import { LoadingSpinner } from "../common/LoadingSpinner";

// PropTypes definitions
const RoomType = PropTypes.shape({
  id: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
  createdAt: PropTypes.string.isRequired,
  participantCount: PropTypes.number,
  waitingCount: PropTypes.number,
  status: PropTypes.number.isRequired,
});

const ErrorMessage = ({ message, onRetry }) => (
  <div className="flex flex-col items-center justify-center h-screen">
    <div className="text-red-500 mb-4">{message}</div>
    <button
      onClick={onRetry}
      className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
    >
      Try Again
    </button>
  </div>
);

ErrorMessage.propTypes = {
  message: PropTypes.string.isRequired,
  onRetry: PropTypes.func.isRequired,
};

const RoomTable = ({ title, rooms, onJoin, joiningRoom }) => (
  <div className="mb-8">
    <div className="overflow-x-auto shadow-lg rounded-lg">
      <table className="min-w-full bg-white border border-gray-200">
        <thead className="bg-gray-200">
          <tr>
            <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">
              Name
            </th>
            <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">
              Created At
            </th>
            <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">
              {title === "Active Meetings"
                ? "Participants"
                : "Total Participants"}
            </th>
            {title === "Active Meetings" && (
              <>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">
                  Waiting
                </th>
                <th className="px-6 py-3 text-center text-sm font-semibold text-gray-700">
                  Action
                </th>
              </>
            )}
          </tr>
        </thead>
        <tbody>
          {rooms.map((room) => (
            <tr key={room.id} className="hover:bg-gray-50 transition-all">
              <td className="border-t border-gray-200 px-6 py-4">
                {room.name}
              </td>
              <td className="border-t border-gray-200 px-6 py-4">
                {formatTimestamp(room.createdAt)}
              </td>
              <td className="border-t border-gray-200 px-6 py-4">
                {room.participantCount || 0}
              </td>
              {title === "Active Meetings" && (
                <>
                  <td className="border-t border-gray-200 px-6 py-4">
                    {room.waitingCount || 0}
                  </td>
                  <td className="border-t border-gray-200 px-6 py-4 text-center">
                    <button
                      onClick={() => onJoin(room)}
                      disabled={joiningRoom === room.id}
                      className="px-4 py-2 rounded-full bg-blue-700 hover:bg-green-500 text-white 
                               transition duration-200 ease-in-out shadow-md transform hover:scale-105 
                               disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      {joiningRoom === room.id ? "Joining..." : "Join"}
                    </button>
                  </td>
                </>
              )}
            </tr>
          ))}
          {rooms.length === 0 && (
            <tr>
              <td
                colSpan={title === "Active Meetings" ? "5" : "3"}
                className="px-6 py-4 text-center text-gray-500"
              >
                No {title.toLowerCase()}
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  </div>
);

RoomTable.propTypes = {
  title: PropTypes.string.isRequired,
  rooms: PropTypes.arrayOf(RoomType).isRequired,
  onJoin: PropTypes.func,
  joiningRoom: PropTypes.string,
};

const Reception = () => {
  const navigate = useNavigate();
  const [roomsWithStatus, setRoomsWithStatus] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [isCreatePanelOpen, setIsCreatePanelOpen] = useState(false);
  const [refreshing, setRefreshing] = useState(false);
  const [joiningRoom, setJoiningRoom] = useState(null);
  const [activeTab, setActiveTab] = useState("active");

  const loadRooms = async () => {
    try {
      const data = await fetchRooms();
      const roomsWithStatusPromises = data.map(async (room) => {
        try {
          const status = await getRoomStatus(room.id);
          return {
            ...room,
            ...status,
          };
        } catch (error) {
          console.error(`Failed to fetch status for room ${room.id}:`, error);
          return room;
        }
      });

      const updatedRooms = await Promise.all(roomsWithStatusPromises);
      setRoomsWithStatus(updatedRooms);
    } catch (error) {
      setError(error.message || "Failed to load rooms");
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  };

  useEffect(() => {
    loadRooms();
  }, [isCreatePanelOpen]);

  const handleJoin = async (room) => {
    try {
      setJoiningRoom(room.id);
      const response = await joinRoom(room.id);
      navigate(`/media/${room.id}`, {
        state: {
          room,
          status: response.status,
        },
      });
    } catch (error) {
      setError("Failed to join room: " + (error.message || "Unknown error"));
    } finally {
      setJoiningRoom(null);
    }
  };

  const handleLogout = async () => {
    try {
      await logout();
      navigate("/auth");
    } catch (error) {
      setError("Failed to logout: " + (error.message || "Unknown error"));
    }
  };

  if (loading) return <LoadingSpinner />;
  if (error) return <ErrorMessage message={error} onRetry={loadRooms} />;

  const activeRooms = roomsWithStatus.filter((room) => room.status === 1);
  const inactiveRooms = roomsWithStatus.filter((room) => room.status === 0);

  return (
    <>
      <Disclosure as="nav" className="bg-black">
        {({ open }) => (
          <>
            <div className="mx-auto max-w-7xl px-2 sm:px-6 lg:px-8">
              <div className="relative flex h-16 items-center justify-between">
                <div className="flex flex-1 items-center justify-center sm:items-stretch sm:justify-start">
                  <div className="flex shrink-0 items-center">
                    <img
                      alt="Chime Cast"
                      src="/image.png"
                      className="h-8 w-auto"
                    />
                  </div>
                  <div className="hidden sm:ml-6 sm:block">
                    <div className="flex space-x-4">
                      <button
                        onClick={() => setActiveTab("active")}
                        className={`px-3 py-2 rounded-md text-sm font-medium ${
                          activeTab === "active"
                            ? "bg-gray-900 text-white"
                            : "text-gray-300 hover:bg-gray-700 hover:text-white"
                        }`}
                      >
                        Active
                      </button>
                      <button
                        onClick={() => setActiveTab("inactive")}
                        className={`px-3 py-2 rounded-md text-sm font-medium ${
                          activeTab === "inactive"
                            ? "bg-gray-900 text-white"
                            : "text-gray-300 hover:bg-gray-700 hover:text-white"
                        }`}
                      >
                        Inactive
                      </button>
                    </div>
                  </div>
                </div>
                <div className="absolute inset-y-0 right-0 flex items-center pr-2 sm:static sm:inset-auto sm:ml-6 sm:pr-0">
                  <button
                    onClick={loadRooms}
                    disabled={refreshing}
                    className="p-2 rounded-full hover:bg-gray-100 transition-all"
                  >
                    <svg
                      className="w-6 h-6 text-white"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
                      />
                    </svg>
                  </button>
                  <Menu>
                    <Menu.Button className="relative flex rounded-full bg-gray-800 text-sm focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-gray-800">
                      <img
                        alt=""
                        src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-1.2.1&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"
                        className="size-8 rounded-full"
                      />
                    </Menu.Button>
                    <Menu.Items className="absolute right-0 z-10 mt-12 w-48 origin-top-right rounded-md bg-white py-1 shadow-lg ring-1 ring-black/5">
                      <Menu.Item>
                        <button
                          onClick={() => setIsCreatePanelOpen(true)}
                          className="block px-4 py-2 text-sm text-gray-700 w-full text-left hover:bg-gray-100"
                        >
                          Create Room
                        </button>
                      </Menu.Item>
                      <Menu.Item>
                        <button
                          onClick={handleLogout}
                          className="block px-4 py-2 text-sm text-gray-700 w-full text-left hover:bg-gray-100"
                        >
                          Sign out
                        </button>
                      </Menu.Item>
                    </Menu.Items>
                  </Menu>
                </div>
              </div>
            </div>
          </>
        )}
      </Disclosure>

      <div className="container mx-auto p-6">
        <CreateRoomPanel
          isOpen={isCreatePanelOpen}
          onClose={() => setIsCreatePanelOpen(false)}
        />

        {activeTab === "active" ? (
          <RoomTable
            title="Active Meetings"
            rooms={activeRooms}
            onJoin={handleJoin}
            joiningRoom={joiningRoom}
          />
        ) : (
          <RoomTable title="Past Meetings" rooms={inactiveRooms} />
        )}
      </div>
    </>
  );
};

export default Reception;
