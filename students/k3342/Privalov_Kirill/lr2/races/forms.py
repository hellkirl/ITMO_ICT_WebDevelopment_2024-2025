from django import forms
from .models import Registration, Comment

class RegistrationForm(forms.ModelForm):
    class Meta:
        model = Registration
        fields = ['race']

    def __init__(self, *args, **kwargs):
        self.user = kwargs.pop('user', None)
        super(RegistrationForm, self).__init__(*args, **kwargs)

    def clean(self):
        cleaned_data = super().clean()
        race = cleaned_data.get('race')
        if race and self.user:
            if Registration.objects.filter(race=race, user=self.user).exists():
                raise forms.ValidationError('Вы уже зарегистрированы на эту гонку.')
        return cleaned_data

class CommentForm(forms.ModelForm):
    class Meta:
        model = Comment
        fields = ['race', 'text', 'comment_type', 'rating']
