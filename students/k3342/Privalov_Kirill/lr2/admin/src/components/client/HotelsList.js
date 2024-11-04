import React, { useEffect, useState } from 'react';
import {
    Container, Grid, Card, CardMedia,
    CardContent, Typography, Button
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import api from '../../services/api';

const HotelsList = () => {
    const [hotels, setHotels] = useState([]);
    const navigate = useNavigate();

    const fetchHotels = async () => {
        try {
            const response = await api.get('/hotels');
            setHotels(response.data);
        } catch (error) {
            alert('Ошибка при получении списка отелей: ' + (error.response?.data || 'Unknown error'));
        }
    };

    useEffect(() => {
        fetchHotels();
    }, []);

    const handleViewDetails = (hotelID) => {
        navigate(`/client/hotels/${hotelID}`);
    };

    return (
        <Container sx={{ mt: 4 }}>
            <Typography variant="h4" gutterBottom>
                Список Отелей
            </Typography>
            <Grid container spacing={4}>
                {hotels.map((hotel) => (
                    <Grid item key={hotel.id} xs={12} sm={6} md={4}>
                        <Card>
                            <CardMedia
                                component="img"
                                height="140"
                                image={hotel.image || 'https://via.placeholder.com/140'}
                                alt={hotel.name}
                            />
                            <CardContent>
                                <Typography variant="h5" component="div">
                                    {hotel.name}
                                </Typography>
                                <Typography variant="body2" color="text.secondary">
                                    {hotel.address}, {hotel.city}, {hotel.country}
                                </Typography>
                                <Button
                                    variant="contained"
                                    color="primary"
                                    sx={{ mt: 2 }}
                                    onClick={() => handleViewDetails(hotel.id)}
                                >
                                    Подробнее
                                </Button>
                            </CardContent>
                        </Card>
                    </Grid>
                ))}
            </Grid>
        </Container>
    );
};

export default HotelsList;