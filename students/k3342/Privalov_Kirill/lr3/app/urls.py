from django.urls import path, include
from rest_framework.routers import DefaultRouter
from .views import (
    PatientViewSet, MedicalCardViewSet, DoctorViewSet, PositionViewSet, LaborContractViewSet,
    ScheduleViewSet, OfficeViewSet, VisitViewSet, DiagnosisViewSet, VisitDiagnosisViewSet,
    ServiceViewSet, ServicePriceViewSet, VisitServiceViewSet, PaymentViewSet
)

router = DefaultRouter()
router.register('patients', PatientViewSet)
router.register('medicalcards', MedicalCardViewSet)
router.register('positions', PositionViewSet)
router.register('doctors', DoctorViewSet)
router.register('laborcontracts', LaborContractViewSet)
router.register('schedules', ScheduleViewSet)
router.register('offices', OfficeViewSet)
router.register('visits', VisitViewSet)
router.register('diagnoses', DiagnosisViewSet)
router.register('visitdiagnoses', VisitDiagnosisViewSet)
router.register('services', ServiceViewSet)
router.register('serviceprices', ServicePriceViewSet)
router.register('visitservices', VisitServiceViewSet)
router.register('payments', PaymentViewSet)

urlpatterns = [
    path('', include(router.urls)),
]