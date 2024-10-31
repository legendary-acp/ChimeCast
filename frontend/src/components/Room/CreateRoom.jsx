// src/components/CreateRoomPanel.jsx
import { useState } from "react";
import { createRoom } from "../../services/roomService";
import PropTypes from "prop-types";

const CreateRoomPanel = ({ isOpen, onClose }) => {
  const [newRoomName, setNewRoomName] = useState("");

  const handleFormSubmit = (e) => {
    e.preventDefault();
    createRoom(newRoomName);
    setNewRoomName("");
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-gray-800 bg-opacity-50 flex justify-center items-center z-50">
      <div className="w-2/5 lg:w-1/2 bg-white p-8 rounded-2xl shadow-2xl">
        <h3 className="text-3xl font-semibold mb-6 text-gray-800">Create New Room</h3>
        <form onSubmit={handleFormSubmit} className="flex flex-col space-y-4">
          <label className="text-lg font-medium text-gray-700">Room Name</label>
          <input
            type="text"
            value={newRoomName}
            onChange={(e) => setNewRoomName(e.target.value)}
            className="p-3 rounded-lg border border-gray-300 focus:outline-none focus:border-blue-500"
            placeholder="Enter room name"
            required
          />
          <button
            type="submit"
            className="bg-blue-600 hover:bg-blue-700 text-white py-3 rounded-lg font-semibold text-lg transition duration-200"
          >
            Create Room
          </button>
          <button
            type="button"
            onClick={onClose}
            className="text-red-500 hover:text-red-600 text-lg mt-4"
          >
            Cancel
          </button>
        </form>
      </div>
    </div>
  );
};

CreateRoomPanel.propTypes = {
  isOpen: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
};

export default CreateRoomPanel;
