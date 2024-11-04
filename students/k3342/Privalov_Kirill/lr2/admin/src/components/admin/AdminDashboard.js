import React, { useEffect, useState } from 'react';
import {
    Container, Typography, Button, Box
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import api from '../../services/api';
import GuestsTable from './GuestsTable';

const AdminDashboard = () => {
    const [guests, setGuests] = useState([]);
    const navigate = useNavigate();

    const fetchGuests = async () => {
        try {
            const response = await api.get('/guests/recent');
            setGuests(response.data);
        } catch (error) {
            alert('Ошибка при получении списка гостей: ' + (error.response?.data || 'Unknown error'));
        }
    };

    useEffect(() => {
        fetchGuests();
    }, []);

    const handleLogout = () => {
        localStorage.removeItem('token');
        navigate('/admin/login');
    };

    const handleCreateReservation = () => {
        navigate('/admin/create-reservation');
    };

    const handleCheckIn = async (reservationID) => {
        try {
            await api.post(`/reservations/${reservationID}/checkin`);
            setGuests((prevGuests) =>
                prevGuests.map((guest) =>
                    guest.reservation_id === reservationID
                        ? { ...guest, status: 'checked_in' }
                        : guest
                )
            );
        } catch (error) {
            alert('Ошибка при заселении: ' + (error.response?.data || 'Unknown error'));
        }
    };

    const handleCheckOut = async (reservationID) => {
        try {
            await api.post(`/reservations/${reservationID}/checkout`);
            setGuests((prevGuests) =>
                prevGuests.map((guest) =>
                    guest.reservation_id === reservationID
                        ? { ...guest, status: 'checked_out' }
                        : guest
                )
            );
        } catch (error) {
            alert('Ошибка при выселении: ' + (error.response?.data || 'Unknown error'));
        }
    };

    return (
        <Container maxWidth="lg">
            <Box sx={{ mt: 4, mb: 2, display: 'flex', justifyContent: 'space-between' }}>
                <Typography variant="h4" component="h1">
                    Административная Панель
                </Typography>
                <Button variant="contained" color="secondary" onClick={handleLogout}>
                    Выйти
                </Button>
            </Box>
            <Button
                variant="contained"
                color="primary"
                onClick={handleCreateReservation}
                sx={{ mb: 2 }}
            >
                Создать Резервацию
            </Button>
            <GuestsTable guests={guests} onCheckIn={handleCheckIn} onCheckOut={handleCheckOut} />
        </Container>
    );
};

export default AdminDashboard;