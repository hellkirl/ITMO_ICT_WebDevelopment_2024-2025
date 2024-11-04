import React, { useEffect, useState } from 'react';
import {
    Container, Typography, Box, Button, Grid, Card, CardContent
} from '@mui/material';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../../services/api';
import Reviews from './Reviews';

const HotelDetails = () => {
    const { id } = useParams();
    const [hotel, setHotel] = useState(null);
    const navigate = useNavigate();

    const fetchHotelDetails = async () => {
        try {
            const response = await api.get(`/hotels/${id}`);
            setHotel(response.data);
        } catch (error) {
            alert('Ошибка при получении деталей отеля: ' + (error.response?.data || 'Unknown error'));
        }
    };

    useEffect(() => {
        fetchHotelDetails();
    }, [id]);

    if (!hotel) {
        return <Typography>Загрузка...</Typography>;
    }

    return (
        <Container sx={{ mt: 4 }}>
            <Button variant="contained" onClick={() => navigate(-1)} sx={{ mb: 2 }}>
                Назад
            </Button>
            <Typography variant="h4" gutterBottom>
                {hotel.name}
            </Typography>
            <Typography variant="subtitle1" gutterBottom>
                {hotel.address}, {hotel.city}, {hotel.country}
            </Typography>
            <Typography variant="body1" gutterBottom>
                Телефон: {hotel.phone} | Email: {hotel.email}
            </Typography>
            <Box sx={{ mt: 4 }}>
                <Typography variant="h5" gutterBottom>
                    Отзывы
                </Typography>
                <Reviews hotelID={id} />
            </Box>
        </Container>
    );
};

export default HotelDetails;