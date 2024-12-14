from rest_framework import serializers
from .models import (
    Patient, MedicalCard, Doctor, Position, LaborContract, Schedule, Office, Visit, Diagnosis,
    VisitDiagnosis, Service, ServicePrice, VisitService, Payment
)

class PatientSerializer(serializers.ModelSerializer):
    class Meta:
        model = Patient
        fields = '__all__'

class MedicalCardSerializer(serializers.ModelSerializer):
    class Meta:
        model = MedicalCard
        fields = '__all__'

class PositionSerializer(serializers.ModelSerializer):
    class Meta:
        model = Position
        fields = '__all__'

class DoctorSerializer(serializers.ModelSerializer):
    class Meta:
        model = Doctor
        fields = '__all__'

class LaborContractSerializer(serializers.ModelSerializer):
    class Meta:
        model = LaborContract
        fields = '__all__'

class ScheduleSerializer(serializers.ModelSerializer):
    class Meta:
        model = Schedule
        fields = '__all__'

class OfficeSerializer(serializers.ModelSerializer):
    class Meta:
        model = Office
        fields = '__all__'

class VisitSerializer(serializers.ModelSerializer):
    class Meta:
        model = Visit
        fields = '__all__'

class DiagnosisSerializer(serializers.ModelSerializer):
    class Meta:
        model = Diagnosis
        fields = '__all__'

class VisitDiagnosisSerializer(serializers.ModelSerializer):
    class Meta:
        model = VisitDiagnosis
        fields = '__all__'

class ServiceSerializer(serializers.ModelSerializer):
    class Meta:
        model = Service
        fields = '__all__'

class ServicePriceSerializer(serializers.ModelSerializer):
    class Meta:
        model = ServicePrice
        fields = '__all__'

class VisitServiceSerializer(serializers.ModelSerializer):
    class Meta:
        model = VisitService
        fields = '__all__'

class PaymentSerializer(serializers.ModelSerializer):
    class Meta:
        model = Payment
        fields = '__all__'
