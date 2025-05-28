<template>
  <div class="container mx-auto p-4">
    <h1 class="text-3xl font-bold mb-4">{{ $t('title') }}</h1>

    <!-- Input Form -->
    <form @submit.prevent="generateGame" class="mb-6">
      <div class="mb-4">
        <label class="block text-sm font-medium">{{ $t('theme') }}</label>
        <input v-model="form.theme" type="text" :placeholder="$t('themePlaceholder')" class="w-full p-2 border rounded" required />
      </div>
      <div class="mb-4">
        <label class="block text-sm font-medium">{{ $t('cardCount') }}</label>
        <input v-model.number="form.cardCount" type="number" min="10" max="100" class="w-full p-2 border rounded" required />
      </div>
      <div class="mb-4">
        <label class="block text-sm font-medium">{{ $t('style') }}</label>
        <select v-model="form.style" class="w-full p-2 border rounded" required>
          <option value="D&D">D&D</option>
          <option value="simple">{{ $t('simple') }}</option>
          <option value="strategy">{{ $t('strategy') }}</option>
        </select>
      </div>
      <div class="mb-4">
        <label class="block text-sm font-medium">{{ $t('description') }}</label>
        <textarea v-model="form.description" :placeholder="$t('descriptionPlaceholder')" class="w-full p-2 border rounded"></textarea>
      </div>
      <button type="submit" class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600">{{ $t('generate') }}</button>
    </form>

    <!-- Saved Games -->
    <div class="mb-6">
      <h2 class="text-2xl font-semibold mb-4">{{ $t('savedGames') }}</h2>
      <div v-if="games.length" class="grid grid-cols-2 gap-4">
        <div v-for="game in games" :key="game.id" class="border p-4 rounded cursor-pointer hover:bg-gray-100" @click="fetchGame(game.id)">
          <h3 class="font-bold">{{ game.theme }}</h3>
          <p>ID: {{ game.id }} | {{ $t('cards') }}: {{ game.card_count }} | {{ $t('style') }}: {{ game.style }}</p>
        </div>
      </div>
      <p v-else>{{ $t('noGames') }}</p>
    </div>

    <!-- Result Preview -->
    <div v-if="selectedGame">
      <h2 class="text-2xl font-semibold mb-4">{{ $t('game') }}: {{ selectedGame.theme }}</h2>
      <p class="mb-4"><strong>{{ $t('story') }}:</strong> {{ selectedGame.description }}</p>
      <h3 class="text-xl font-semibold mb-2">{{ $t('cards') }}</h3>
      <!-- Card Filter -->
      <div class="mb-4">
        <label class="block text-sm font-medium">{{ $t('filterByType') }}</label>
        <select v-model="cardFilter" class="w-full p-2 border rounded">
          <option value="all">{{ $t('all') }}</option>
          <option value="role">{{ $t('role') }}</option>
          <option value="event">{{ $t('event') }}</option>
          <option value="item">{{ $t('item') }}</option>
        </select>
      </div>
      <div class="grid grid-cols-2 gap-4">
        <div v-for="card in filteredCards" :key="card.id" class="border p-4 rounded">
          <h4 class="font-bold">{{ card.name }} ({{ $t(card.type) }})</h4>
          <p>{{ card.description }}</p>
          <p><strong>{{ $t('effect') }}:</strong> {{ card.effect }}</p>
        </div>
      </div>
      <a :href="pdfUrl" class="mt-4 inline-block bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600" target="_blank">{{ $t('downloadPDF') }}</a>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      form: {
        theme: 'Fantasy',
        cardCount: 20,
        style: 'D&D',
        description: '',
      },
      games: [],
      selectedGame: null,
      pdfUrl: '',
      cardFilter: 'all', // New: Card type filter
    };
  },
  computed: {
    filteredCards() {
      if (this.cardFilter === 'all' || !this.selectedGame) {
        return this.selectedGame?.cards || [];
      }
      return this.selectedGame.cards.filter(card => card.type === this.cardFilter);
    },
  },
  async mounted() {
    await this.fetchGames();
  },
  methods: {
    async generateGame() {
      try {
        const response = await fetch('http://localhost:8080/api/v1/game', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(this.form),
        });
        if (!response.ok) {
          throw new Error((await response.json()).error);
        }
        const data = await response.json();
        alert(this.$t('gameGenerated', { id: data.game_id }));
        await this.fetchGames();
        await this.fetchGame(data.game_id);
      } catch (error) {
        console.error('Generation failed:', error);
        alert(this.$t('generateFailed', { error: error.message }));
      }
    },
    async fetchGames() {
      try {
        const response = await fetch('http://localhost:8080/api/v1/games');
        if (!response.ok) {
          throw new Error('Failed to fetch games');
        }
        this.games = await response.json();
      } catch (error) {
        console.error('Fetch games failed:', error);
        alert(this.$t('fetchGamesFailed'));
      }
    },
    async fetchGame(id) {
      try {
        const response = await fetch(`http://localhost:8080/api/v1/games/${id}`);
        if (!response.ok) {
          throw new Error((await response.json()).error);
        }
        this.selectedGame = await response.json();
        this.pdfUrl = `http://localhost:8080/api/v1/generate-pdf/${id}`;
      } catch (error) {
        console.error('Fetch game failed:', error);
        alert(this.$t('fetchGameFailed', { error: error.message }));
      }
    },
  },
};
</script>

<style>
@import 'tailwindcss/tailwind.css';
</style>