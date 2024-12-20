# Generated by Django 5.1.4 on 2024-12-05 17:09

import django.db.models.deletion
from django.conf import settings
from django.db import migrations, models


class Migration(migrations.Migration):

    initial = True

    dependencies = [
        migrations.swappable_dependency(settings.AUTH_USER_MODEL),
    ]

    operations = [
        migrations.CreateModel(
            name='Participant',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('full_name', models.CharField(max_length=255, verbose_name='ФИО участника')),
                ('team_name', models.CharField(max_length=255, verbose_name='Название команды')),
                ('car_description', models.TextField(verbose_name='Описание автомобиля')),
                ('participant_description', models.TextField(verbose_name='Описание участника')),
                ('experience', models.PositiveIntegerField(verbose_name='Опыт (в годах)')),
                ('participant_class', models.CharField(max_length=100, verbose_name='Класс участника')),
            ],
        ),
        migrations.CreateModel(
            name='Race',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('name', models.CharField(max_length=255, verbose_name='Название гонки')),
                ('date', models.DateField(verbose_name='Дата гонки')),
                ('time', models.TimeField(blank=True, null=True, verbose_name='Время заезда')),
                ('result', models.CharField(blank=True, max_length=255, null=True, verbose_name='Результат')),
            ],
        ),
        migrations.CreateModel(
            name='Comment',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('date', models.DateField(verbose_name='Дата заезда')),
                ('text', models.TextField(verbose_name='Текст комментария')),
                ('comment_type', models.CharField(choices=[('collaboration', 'Вопрос о сотрудничестве'), ('race', 'Вопрос о гонках'), ('other', 'Иное')], max_length=20, verbose_name='Тип комментария')),
                ('rating', models.PositiveIntegerField(default=1, verbose_name='Рейтинг')),
                ('user', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to=settings.AUTH_USER_MODEL, verbose_name='Комментатор')),
                ('race', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='races.race', verbose_name='Гонка')),
            ],
        ),
        migrations.CreateModel(
            name='Registration',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('participant', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='races.participant', verbose_name='Участник')),
                ('race', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='races.race', verbose_name='Гонка')),
                ('user', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to=settings.AUTH_USER_MODEL, verbose_name='Пользователь')),
            ],
        ),
    ]
