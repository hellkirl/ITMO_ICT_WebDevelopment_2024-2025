import React from 'react';
import { Navigate } from 'react-router-dom';
import jwtDecode from 'jwt-decode';

const ProtectedRoute = ({ children, role }) => {
    const token = localStorage.getItem('token');
    if (!token) {
        return <Navigate to="/admin/login" replace />;
    }

    try {
        const decoded = jwtDecode(token);
        if (role && decoded.role !== role) {
            return <Navigate to="/admin/login" replace />;
        }
    } catch (error) {
        return <Navigate to="/admin/login" replace />;
    }

    return children;
};

export default ProtectedRoute;