// src/routes.js
import ToggleAuth from "./components/Auth/ToggleAuth";
import { Navigate } from 'react-router-dom';
import RoomList from "./components/Room/RoomList";
import MediaPreview from "./components/Room/MediaPreview";

const routes = [
  { path: '/auth', element: <ToggleAuth /> },
  { path: '/', element: <Navigate to="/auth" /> },
  { path: '/rooms', element: <RoomList/> },
  { path: '/media/:roomId', element: <MediaPreview /> }
  // Add more routes here as needed
];

export default routes;
