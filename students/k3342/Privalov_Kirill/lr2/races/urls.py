from django.urls import path
from . import views

urlpatterns = [
    path('', views.race_list, name='race_list'),
    path('register/', views.register_racer, name='register_racer'),
    path('race/<int:pk>/', views.race_detail, name='race_detail'),
    path('race/<int:pk>/delete/', views.delete_registration, name='delete_registration'),
    path('comment/add/', views.add_comment, name='add_comment'),
]