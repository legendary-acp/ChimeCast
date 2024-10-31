// src/routes.js
import ToggleAuth from "./components/Auth/ToggleAuth";
import { Navigate } from 'react-router-dom';
import RoomList from "./components/Room/RoomList";

const routes = [
  { path: '/auth', element: <ToggleAuth /> },
  { path: '/', element: <Navigate to="/auth" /> },
  { path: '/rooms', element: <RoomList/> }
  // Add more routes here as needed
];

export default routes;
