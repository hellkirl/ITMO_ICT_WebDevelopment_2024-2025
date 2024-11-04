import React, { useState, useEffect } from 'react';
import {
    Container, TextField, Button, Typography, Box, MenuItem
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import api from '../../services/api';

const CreateReservation = () => {
    const [users, setUsers] = useState([]);
    const [rooms, setRooms] = useState([]);
    const [userID, setUserID] = useState('');
    const [roomID, setRoomID] = useState('');
    const [checkIn, setCheckIn] = useState('');
    const [checkOut, setCheckOut] = useState('');
    const navigate = useNavigate();

    const fetchUsers = async () => {
        try {
            const response = await api.get('/users');
            setUsers(response.data);
        } catch (error) {
            alert('Ошибка при получении пользователей: ' + (error.response?.data || 'Unknown error'));
        }
    };

    const fetchRooms = async () => {
        try {
            const response = await api.get('/rooms');
            setRooms(response.data);
        } catch (error) {
            alert('Ошибка при получении комнат: ' + (error.response?.data || 'Unknown error'));
        }
    };

    useEffect(() => {
        fetchUsers();
        fetchRooms();
    }, []);

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            await api.post('/reservations', {
                user_id: userID,
                room_id: roomID,
                check_in: checkIn,
                check_out: checkOut,
            });
            alert('Резервация успешно создана');
            navigate('/admin/dashboard');
        } catch (error) {
            alert('Ошибка при создании резервации: ' + (error.response?.data || 'Unknown error'));
        }
    };

    return (
        <Container maxWidth="sm">
            <Box sx={{ mt: 4 }}>
                <Typography variant="h5" component="h2" gutterBottom>
                    Создать Резервацию
                </Typography>
                <form onSubmit={handleSubmit}>
                    <TextField
                        select
                        label="Пользователь"
                        fullWidth
                        required
                        margin="normal"
                        value={userID}
                        onChange={(e) => setUserID(e.target.value)}
                    >
                        {users.map((user) => (
                            <MenuItem key={user.id} value={user.id}>
                                {user.username} ({user.email})
                            </MenuItem>
                        ))}
                    </TextField>
                    <TextField
                        select
                        label="Комната"
                        fullWidth
                        required
                        margin="normal"
                        value={roomID}
                        onChange={(e) => setRoomID(e.target.value)}
                    >
                        {rooms.map((room) => (
                            <MenuItem key={room.id} value={room.id}>
                                {room.number} - {room.type} (${room.price}/ночь)
                            </MenuItem>
                        ))}
                    </TextField>
                    <TextField
                        label="Дата Заезда"
                        type="date"
                        fullWidth
                        required
                        margin="normal"
                        InputLabelProps={{
                            shrink: true,
                        }}
                        value={checkIn}
                        onChange={(e) => setCheckIn(e.target.value)}
                    />
                    <TextField
                        label="Дата Выезда"
                        type="date"
                        fullWidth
                        required
                        margin="normal"
                        InputLabelProps={{
                            shrink: true,
                        }}
                        value={checkOut}
                        onChange={(e) => setCheckOut(e.target.value)}
                    />
                    <Button
                        type="submit"
                        variant="contained"
                        color="primary"
                        fullWidth
                        sx={{ mt: 2 }}
                    >
                        Создать
                    </Button>
                </form>
            </Box>
        </Container>
    );
};

export default CreateReservation;