from django.urls import path

from .views import (
    WarriorAPIView,
    WarriorDeleteView,
    WarriorUpdateView,
    OccupationCreateView,
    SkillAPIView,
    SkillCreateView,
    WarriorOccupationAPIView,
    WarriorSkillAPIView,
    WarriorSkillOccupationAPIView
)

urlpatterns = [
    path('warriors/', WarriorAPIView.as_view(), name='warrior-list'),
    path('warrior/delete/<int:pk>/', WarriorDeleteView.as_view(), name='warrior-delete'),
    path('warrior/update/<int:pk>/', WarriorUpdateView.as_view(), name='warrior-update'),
    path('occupation/create/', OccupationCreateView.as_view(), name='occupation-create'),
    path('skills/', SkillAPIView.as_view(), name='skill-list'),
    path('skill/create/', SkillCreateView.as_view(), name='skill-create'),
    path('warriors_and_occupations/', WarriorOccupationAPIView.as_view(), name='warrior-occupations'),
    path('warriors_and_skills/', WarriorSkillAPIView.as_view(), name='warrior-skills'),
    path('warriors_and_skills_and_occupations/', WarriorSkillOccupationAPIView.as_view(),
         name='warrior-skills-occupations'),
]
