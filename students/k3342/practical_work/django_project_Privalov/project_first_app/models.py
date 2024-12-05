from django.conf import settings
from django.contrib.auth.models import AbstractUser
from django.db import models

class Car(models.Model):
    license_plate = models.CharField(max_length=15, verbose_name="License Plate", unique=True)
    brand = models.CharField(max_length=20, verbose_name="Brand")
    model = models.CharField(max_length=20, verbose_name="Model")
    color = models.CharField(max_length=30, null=True, blank=True, verbose_name="Color")

    def __str__(self):
        return f"{self.brand} {self.model} ({self.license_plate})"

    class Meta:
        verbose_name = "Car"
        verbose_name_plural = "Cars"
        ordering = ['brand', 'model']


class Owner(AbstractUser):
    birth_date = models.DateField(verbose_name="Birth Date", null=True, blank=True)
    passport_number = models.CharField(max_length=10, verbose_name="Passport Number", unique=True)
    home_address = models.CharField(max_length=255, verbose_name="Home Address")
    nationality = models.CharField(max_length=50, verbose_name="Nationality")
    cars = models.ManyToManyField(Car, through='Ownership', related_name='owners')

    groups = models.ManyToManyField(
        'auth.Group',
        related_name='project_first_app_owner_set',
        blank=True,
        help_text='The groups this user belongs to. A user will get all permissions granted to each of their groups.',
        related_query_name='user',
    )
    user_permissions = models.ManyToManyField(
        'auth.Permission',
        related_name='project_first_app_owner_set',
        blank=True,
        help_text='Specific permissions for this user.',
        related_query_name='user',
    )

    def __str__(self):
        return f"{self.first_name} {self.last_name}"


class Ownership(models.Model):
    owner = models.ForeignKey(settings.AUTH_USER_MODEL, on_delete=models.CASCADE, verbose_name="Owner")
    car = models.ForeignKey(Car, on_delete=models.CASCADE, verbose_name="Car")
    start_date = models.DateField(verbose_name="Start Date")
    end_date = models.DateField(null=True, blank=True, verbose_name="End Date")

    def __str__(self):
        return f"{self.owner} owns {self.car} from {self.start_date} to {self.end_date or 'Present'}"

    class Meta:
        verbose_name = "Ownership"
        verbose_name_plural = "Ownerships"
        unique_together = ('owner', 'car', 'start_date')


class DriversLicense(models.Model):
    LICENSE_TYPES = [
        ('A', 'Type A'),
        ('B', 'Type B'),
        ('C', 'Type C'),
    ]

    owner = models.ForeignKey(settings.AUTH_USER_MODEL, on_delete=models.CASCADE, related_name='drivers_licenses', verbose_name="Owner")
    license_number = models.CharField(max_length=10, verbose_name="License Number", unique=True)
    license_type = models.CharField(max_length=2, choices=LICENSE_TYPES, verbose_name="License Type")
    issue_date = models.DateField(verbose_name="Issue Date")

    def __str__(self):
        return f"{self.license_type} License {self.license_number} for {self.owner}"

    class Meta:
        verbose_name = "Driver's License"
        verbose_name_plural = "Driver's Licenses"
        ordering = ['-issue_date']
