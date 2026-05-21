<script setup>
import { ArrowRight, Star, ShieldCheck, Globe, Clock } from 'lucide-vue-next';
</script>

<template>
  <div class="space-y-24 pb-20">
    <!-- Hero Section -->
    <section class="relative overflow-hidden rounded-[3rem] bg-slate-900 py-24 text-white lg:py-32">
      <div class="absolute inset-0 z-0 opacity-40">
         <img 
           src="https://images.unsplash.com/photo-1542314831-068cd1dbfeeb?auto=format&fit=crop&w=1920&q=80" 
           class="h-full w-full object-cover" 
           alt="Luxury Hotel"
         />
         <div class="absolute inset-0 bg-gradient-to-t from-slate-900 via-slate-900/60 to-transparent"></div>
      </div>
      
      <div class="relative z-10 mx-auto max-w-4xl px-6 text-center">
        <div>
          <span class="mb-6 inline-block rounded-full bg-blue-600/20 px-4 py-2 text-sm font-bold tracking-widest text-blue-400 backdrop-blur-sm">
            WELCOME TO LUMINA
          </span>
          <h1 class="mb-8 text-5xl font-extrabold tracking-tight sm:text-7xl">
            Sanctuary of <span class="text-blue-500">Luxury</span> <br class="hidden sm:block" /> and Serenity.
          </h1>
          <p class="mb-10 text-xl text-slate-300">
            Discover a world of refined elegance where every detail is designed to provide an unforgettable experience.
          </p>
          <div class="flex flex-col items-center justify-center gap-4 sm:flex-row">
            <NuxtLink to="/rooms" class="group flex items-center gap-2 rounded-full bg-blue-600 px-8 py-4 text-lg font-bold text-white shadow-xl shadow-blue-500/20 transition-all hover:bg-blue-700 hover:scale-105 active:scale-95">
              Explore Our Rooms
              <ArrowRight class="h-5 w-5 transition-transform group-hover:translate-x-1" />
            </NuxtLink>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section>
      <div class="mb-16 text-center">
         <h2 class="text-4xl font-extrabold tracking-tight text-slate-900">Experience Excellence</h2>
         <p class="mt-4 text-xl text-slate-500">Tailored services for your ultimate comfort.</p>
      </div>
      
      <div class="grid gap-12 sm:grid-cols-2 lg:grid-cols-3">
        <div v-for="(feature, idx) in [
          { icon: Globe, title: 'Prime Locations', desc: 'Centrally located in the most vibrant cities.' },
          { icon: ShieldCheck, title: 'Verified Luxury', desc: 'Rooms inspected for the highest standards.' },
          { icon: Clock, title: '24/7 Concierge', desc: 'Our team is ready to assist you anytime.' },
        ]" :key="idx" class="flex flex-col items-center text-center">
          <div class="mb-6 flex h-16 w-16 items-center justify-center rounded-2xl bg-blue-50 text-blue-600">
            <component :is="feature.icon" class="h-8 w-8" />
          </div>
          <h3 class="mb-3 text-xl font-bold text-slate-900">{{ feature.title }}</h3>
          <p class="text-slate-500">{{ feature.desc }}</p>
        </div>
      </div>
    </section>
  </div>
</template>
