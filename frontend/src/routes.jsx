// src/routes.js
import ToggleAuth from "./components/Auth/ToggleAuth";
import { Navigate } from 'react-router-dom';
import {
  RoomList,
  MediaPreview,
  VideoCall,
} from "./components/Room";


const routes = [
  { path: '/auth', element: <ToggleAuth /> },
  { path: '/', element: <Navigate to="/auth" /> },
  { path: '/rooms', element: <RoomList/> },
  { path: '/media/:roomId', element: <MediaPreview /> },
  { path: '/meeting/:roomId', element: <VideoCall /> }
  // Add more routes here as needed
];

export default routes;
