import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
    Container, TextField, Button, Typography, Box
} from '@mui/material';
import api from '../../services/api';

const AdminLogin = () => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const response = await api.post('/users/login', { email, password });
            localStorage.setItem('token', response.data.token);
            navigate('/admin/dashboard');
        } catch (error) {
            alert('Ошибка при входе: ' + (error.response?.data || 'Unknown error'));
        }
    };

    return (
        <Container maxWidth="sm">
            <Box sx={{ mt: 8 }}>
                <Typography variant="h4" component="h1" gutterBottom>
                    Администратор Вход
                </Typography>
                <form onSubmit={handleSubmit}>
                    <TextField
                        label="Email"
                        type="email"
                        fullWidth
                        required
                        margin="normal"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                    />
                    <TextField
                        label="Пароль"
                        type="password"
                        fullWidth
                        required
                        margin="normal"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                    />
                    <Button
                        type="submit"
                        variant="contained"
                        color="primary"
                        fullWidth
                        sx={{ mt: 2 }}
                    >
                        Войти
                    </Button>
                </form>
            </Box>
        </Container>
    );
};

export default AdminLogin;