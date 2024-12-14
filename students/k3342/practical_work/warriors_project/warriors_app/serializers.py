from rest_framework import serializers

from .models import Warrior, Skill, Occupation


class WarriorSerializer(serializers.ModelSerializer):
    class Meta:
        model = Warrior
        fields = "__all__"


class OccupationSerializer(serializers.ModelSerializer):
    class Meta:
        model = Occupation
        fields = "__all__"


class OccupationCreateSerializer(serializers.ModelSerializer):
    class Meta:
        model = Occupation
        fields = "__all__"


class SkillSerializer(serializers.ModelSerializer):
    class Meta:
        model = Skill
        fields = "__all__"


class WarriorOccupationSerializer(serializers.ModelSerializer):
    occupation = serializers.SlugRelatedField(
        read_only=True, slug_field='title'
    )

    class Meta:
        model = Warrior
        fields = "__all__"


class WarriorRelatedSerializer(serializers.ModelSerializer):
    skill = serializers.SlugRelatedField(
        read_only=True, many=True, slug_field='title'
    )

    class Meta:
        model = Warrior
        fields = "__all__"


class WarriorSkillSerializer(serializers.ModelSerializer):
    skill = SkillSerializer(many=True, read_only=True)

    class Meta:
        model = Warrior
        fields = "__all__"


class WarriorNestedSerializer(serializers.ModelSerializer):
    occupation = OccupationSerializer(read_only=True)
    skill = SkillSerializer(many=True, read_only=True)
    race = serializers.CharField(source='get_race_display', read_only=True)

    class Meta:
        model = Warrior
        fields = "__all__"
