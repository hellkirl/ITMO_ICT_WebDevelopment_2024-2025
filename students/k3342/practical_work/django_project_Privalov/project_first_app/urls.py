from django.urls import path
from .views import (
    CarListView,
    CarDetailView,
    CarCreateView,
    CarUpdateView,
    CarDeleteView,
    OwnerListView,
    OwnerDetailView,
    OwnerCreateView,
    OwnerUpdateView,
    OwnerDeleteView,
)

urlpatterns = [
    path('owners/', OwnerListView.as_view(), name='owner-list'),
    path('owners/<int:pk>/', OwnerDetailView.as_view(), name='owner-detail'),
    path('owners/create/', OwnerCreateView.as_view(), name='owner-create'),
    path('owners/<int:pk>/update/', OwnerUpdateView.as_view(), name='owner-update'),
    path('owners/<int:pk>/delete/', OwnerDeleteView.as_view(), name='owner-delete'),

    path('cars/', CarListView.as_view(), name='car-list'),
    path('cars/<int:pk>/', CarDetailView.as_view(), name='car-detail'),
    path('cars/create/', CarCreateView.as_view(), name='car-create'),
    path('cars/<int:pk>/update/', CarUpdateView.as_view(), name='car-update'),
    path('cars/<int:pk>/delete/', CarDeleteView.as_view(), name='car-delete'),
]
