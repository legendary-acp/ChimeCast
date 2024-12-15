// src/routes.js
import ToggleAuth from "./components/Auth/ToggleAuth";
import AuthRedirect from "./components/Auth/AuthRedirect";
import { MediaPreview, VideoCall } from "./components/Room";
import ProtectedRoute from "./components/common/ProtectedRoute";
import Reception from "./components/Reception/Reception";

const routes = [
  { 
    path: '/auth', 
    element: <ToggleAuth /> 
  },
  { 
    path: '/', 
    element: <AuthRedirect /> 
  },
  { 
    path: '/rooms', 
    element: <ProtectedRoute><Reception /></ProtectedRoute> 
  },
  { 
    path: '/media/:roomId', 
    element: <ProtectedRoute><MediaPreview /></ProtectedRoute> 
  },
  { 
    path: '/meeting/:roomId', 
    element: <ProtectedRoute><VideoCall /></ProtectedRoute> 
  }
];

export default routes;