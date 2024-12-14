from django.db import models


class Warrior(models.Model):
    races = (
        ('j', 'junior'),
        ('m', 'middle'),
        ('s', 'senior'),
    )
    race = models.CharField(max_length=1, choices=races, verbose_name='Расса')
    name = models.CharField(max_length=120, verbose_name='Имя')
    level = models.IntegerField(verbose_name='Уровень', default=0)
    skill = models.ManyToManyField(
        'Skill',
        verbose_name='Навыки',
        through='SkillOfWarrior',
        related_name='warrior_skills'
    )
    profession = models.ForeignKey(
        'Occupation',
        on_delete=models.CASCADE,
        verbose_name='Вид деятельности',
        blank=True,
        null=True
    )

    def __str__(self):
        return "{} {} {} {}".format(self.race, self.name, self.level, self.profession)


class Skill(models.Model):
    title = models.CharField(max_length=120, verbose_name='Наименование')

    def __str__(self):
        return self.title


class SkillOfWarrior(models.Model):
    skill = models.ForeignKey('Skill', verbose_name='Умение', on_delete=models.CASCADE)
    warrior = models.ForeignKey('Warrior', verbose_name='Воин', on_delete=models.CASCADE)
    level = models.IntegerField(verbose_name='Уровень освоения умения')

    def __str__(self):
        return "{} has {} skill".format(self.warrior, self.skill)


class Occupation(models.Model):
    title = models.CharField(max_length=120, verbose_name='Название')
    description = models.TextField(verbose_name='Описание')

    def __str__(self):
        return self.title
