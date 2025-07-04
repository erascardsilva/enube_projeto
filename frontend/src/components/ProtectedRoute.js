import { Navigate } from "react-router-dom";
import PropTypes from "prop-types";
import authService from "services/authService";

const ProtectedRoute = ({ children }) => {
  if (!authService.isAuthenticated()) {
    return (
      <Navigate to="/authentication/sign-in" replace />
    );
  }
  return children;
};

ProtectedRoute.propTypes = {
  children: PropTypes.node.isRequired,
};

export default ProtectedRoute; 