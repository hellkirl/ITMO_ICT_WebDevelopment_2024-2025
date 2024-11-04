import React, { useEffect, useState } from 'react';
import {
    List, ListItem, ListItemText, Rating, Typography, Box
} from '@mui/material';
import api from '../../services/api';

const Reviews = ({ hotelID }) => {
    const [reviews, setReviews] = useState([]);

    const fetchReviews = async () => {
        try {
            const response = await api.get(`/hotels/${hotelID}/reviews`);
            setReviews(response.data);
        } catch (error) {
            alert('Ошибка при получении отзывов: ' + (error.response?.data || 'Unknown error'));
        }
    };

    useEffect(() => {
        fetchReviews();
    }, [hotelID]);

    if (reviews.length === 0) {
        return <Typography>Нет отзывов.</Typography>;
    }

    return (
        <List>
            {reviews.map((review) => (
                <ListItem key={review.id} alignItems="flex-start">
                    <ListItemText
                        primary={
                            <Box sx={{ display: 'flex', alignItems: 'center' }}>
                                <Rating value={review.rating} readOnly />
                                <Typography variant="subtitle2" sx={{ ml: 1 }}>
                                    {review.user.username}
                                </Typography>
                            </Box>
                        }
                        secondary={review.comment}
                    />
                </ListItem>
            ))}
        </List>
    );
};

export default Reviews;