// src/components/Auth/ProtectedRoute.jsx
import { useEffect, useState } from "react";
import { Navigate } from "react-router-dom";
import { validateAuth } from "../../services/authService";
import PropTypes from "prop-types"; // Add this import

const ProtectedRoute = ({ children }) => {
  const [authState, setAuthState] = useState({
    isChecking: true,
    isAuthenticated: false,
  });

  useEffect(() => {
    const checkAuth = async () => {
      const result = await validateAuth();
      setAuthState({
        isChecking: false,
        isAuthenticated: !!result,
      });
    };

    checkAuth();
  }, []);

  if (authState.isChecking) {
    return <div>Loading...</div>;
  }

  return authState.isAuthenticated ? children : <Navigate to="/auth" replace />;
};

// Add PropTypes validation
ProtectedRoute.propTypes = {
  children: PropTypes.node.isRequired,
};

export default ProtectedRoute;
