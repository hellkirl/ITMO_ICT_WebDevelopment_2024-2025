from rest_framework import generics
from rest_framework.response import Response
from rest_framework.views import APIView

from .models import Warrior, Skill
from .serializers import (
    WarriorSerializer,
    OccupationCreateSerializer,
    SkillSerializer,
    WarriorOccupationSerializer,
    WarriorSkillSerializer,
    WarriorNestedSerializer
)


class WarriorAPIView(APIView):
    def get(self, request):
        warriors = Warrior.objects.all()
        serializer = WarriorSerializer(warriors, many=True)
        return Response({"Warriors": serializer.data})


class WarriorDeleteView(generics.DestroyAPIView):
    serializer_class = WarriorSerializer
    queryset = Warrior.objects.all()
    lookup_field = 'id'


class WarriorUpdateView(generics.UpdateAPIView):
    serializer_class = WarriorSerializer
    queryset = Warrior.objects.all()


class WarriorOccupationAPIView(APIView):
    def get(self, request):
        warriors = Warrior.objects.select_related('profession').all()
        serializer = WarriorOccupationSerializer(warriors, many=True)
        return Response({"Warriors": serializer.data})


class WarriorSkillAPIView(APIView):
    def get(self, request):
        warriors = Warrior.objects.prefetch_related('skill').all()
        serializer = WarriorSkillSerializer(warriors, many=True)
        return Response({"Warriors": serializer.data})


class WarriorSkillOccupationAPIView(APIView):
    def get(self, request):
        warriors = Warrior.objects.select_related('profession').prefetch_related('skill').all()
        serializer = WarriorNestedSerializer(warriors, many=True)
        return Response({"Warriors": serializer.data})


class OccupationCreateView(APIView):
    def post(self, request):
        occupation_data = request.data.get("occupation")
        serializer = OccupationCreateSerializer(data=occupation_data)

        if serializer.is_valid(raise_exception=True):
            occupation_saved = serializer.save()

        return Response({"Success": f"Occupation '{occupation_saved.title}' created successfully."})


class SkillAPIView(APIView):
    def get(self, request):
        skills = Skill.objects.all()
        serializer = SkillSerializer(skills, many=True)
        return Response({"Skills": serializer.data})


class SkillCreateView(APIView):
    def post(self, request):
        skill_data = request.data.get("skill")
        serializer = SkillSerializer(data=skill_data)

        if serializer.is_valid(raise_exception=True):
            skill_saved = serializer.save()

        return Response({"Success": f"Skill '{skill_saved.title}' created successfully."})
