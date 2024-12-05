from django.http import Http404
from django.shortcuts import render, redirect
from django.urls import reverse_lazy
from django.views import View
from django.views.generic import (
    ListView,
    DetailView,
    CreateView,
    UpdateView,
    DeleteView,
)

from .forms import OwnerForm
from .models import Car, Owner


class CarListView(ListView):
    model = Car
    context_object_name = 'cars'
    template_name = 'cars/list.html'


class CarDetailView(DetailView):
    model = Car
    context_object_name = 'car'
    template_name = 'cars/detail.html'


class CarCreateView(CreateView):
    model = Car
    fields = '__all__'
    template_name = 'cars/form.html'
    success_url = reverse_lazy('car-list')


class CarUpdateView(UpdateView):
    model = Car
    fields = '__all__'
    template_name = 'cars/form.html'
    success_url = reverse_lazy('car-list')


class CarDeleteView(DeleteView):
    model = Car
    template_name = 'cars/confirm_delete.html'
    success_url = reverse_lazy('car-list')


class CarBulkDeleteView(View):
    template_name = 'cars/bulk_confirm_delete.html'
    success_url = reverse_lazy('car-list')

    def post(self, request, *args, **kwargs):
        car_ids = request.POST.getlist('selected_cars')
        Car.objects.filter(id__in=car_ids).delete()
        return redirect(self.success_url)

    def get(self, request, *args, **kwargs):
        cars = Car.objects.all()
        return render(request, self.template_name, {'cars': cars})


class OwnerListView(ListView):
    model = Owner
    context_object_name = 'owners'
    template_name = 'owners/list.html'


class OwnerDetailView(DetailView):
    model = Owner
    context_object_name = 'owner'
    template_name = 'owners/detail.html'


class OwnerCreateView(CreateView):
    model = Owner
    form_class = OwnerForm
    template_name = 'owners/form.html'
    success_url = reverse_lazy('owner-list')


class OwnerUpdateView(UpdateView):
    model = Owner
    form_class = OwnerForm
    template_name = 'owners/form.html'
    success_url = reverse_lazy('owner-list')


class OwnerDeleteView(DeleteView):
    model = Owner
    template_name = 'owners/confirm_delete.html'
    success_url = reverse_lazy('owner-list')


class OwnerBulkDeleteView(View):
    template_name = 'owners/bulk_confirm_delete.html'
    success_url = reverse_lazy('owner-list')

    def post(self, request, *args, **kwargs):
        owner_ids = request.POST.getlist('selected_owners')
        Owner.objects.filter(id__in=owner_ids).delete()
        return redirect(self.success_url)

    def get(self, request, *args, **kwargs):
        owners = Owner.objects.all()
        return render(request, self.template_name, {'owners': owners})
