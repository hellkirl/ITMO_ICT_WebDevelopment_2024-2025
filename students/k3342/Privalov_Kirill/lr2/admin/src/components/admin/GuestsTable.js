import React from 'react';
import {
    Table, TableBody, TableCell, TableContainer,
    TableHead, TableRow, Paper, Button
} from '@mui/material';

const GuestsTable = ({ guests, onCheckIn, onCheckOut }) => {
    return (
        <TableContainer component={Paper}>
            <Table>
                <TableHead>
                    <TableRow>
                        <TableCell>Reservation ID</TableCell>
                        <TableCell>User ID</TableCell>
                        <TableCell>Имя пользователя</TableCell>
                        <TableCell>Номер комнаты</TableCell>
                        <TableCell>Заезд</TableCell>
                        <TableCell>Выезд</TableCell>
                        <TableCell>Статус</TableCell>
                        <TableCell>Действия</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {guests.map((guest) => (
                        <TableRow key={guest.reservation_id}>
                            <TableCell>{guest.reservation_id}</TableCell>
                            <TableCell>{guest.user_id}</TableCell>
                            <TableCell>{guest.username}</TableCell>
                            <TableCell>{guest.room_number}</TableCell>
                            <TableCell>{new Date(guest.check_in).toLocaleDateString()}</TableCell>
                            <TableCell>{new Date(guest.check_out).toLocaleDateString()}</TableCell>
                            <TableCell>{guest.status}</TableCell>
                            <TableCell>
                                {guest.status === 'reserved' && (
                                    <Button
                                        variant="contained"
                                        color="primary"
                                        onClick={() => onCheckIn(guest.reservation_id)}
                                        sx={{ mr: 1 }}
                                    >
                                        Заселить
                                    </Button>
                                )}
                                {guest.status === 'checked_in' && (
                                    <Button
                                        variant="contained"
                                        color="secondary"
                                        onClick={() => onCheckOut(guest.reservation_id)}
                                    >
                                        Выселить
                                    </Button>
                                )}
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </TableContainer>
    );
};

export default GuestsTable;