<template>
  <div class="container mx-auto p-4">
    <h1 class="text-3xl font-bold mb-4">虛擬桌遊生成器</h1>
    
    <!-- 輸入表單 -->
    <form @submit.prevent="generateGame" class="mb-6">
      <div class="mb-4">
        <label class="block text-sm font-medium">桌遊主題</label>
        <input v-model="form.theme" type="text" placeholder="e.g., 奇幻冒險" class="w-full p-2 border rounded" required />
      </div>
      <div class="mb-4">
        <label class="block text-sm font-medium">卡牌數量</label>
        <input v-model.number="form.cardCount" type="number" min="1" max="20" class="w-full p-2 border rounded" required />
      </div>
      <div class="mb-4">
        <label class="block text-sm font-medium">遊戲風格</label>
        <select v-model="form.style" class="w-full p-2 border rounded" required>
          <option value="simple">簡單</option>
          <option value="strategy">策略</option>
        </select>
      </div>
      <button type="submit" class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600">生成桌遊</button>
    </form>

    <!-- 結果預覽 -->
    <div v-if="cards.length">
      <h2 class="text-2xl font-semibold mb-4">生成結果</h2>
      <div class="grid grid-cols-2 gap-4">
        <div v-for="card in cards" :key="card.id" class="border p-4 rounded">
          <h3 class="font-bold">{{ card.name }}</h3>
          <p>{{ card.description }}</p>
          <p><strong>效果：</strong> {{ card.effect }}</p>
        </div>
      </div>
      <a :href="pdfUrl" class="mt-4 inline-block bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600">下載PDF</a>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      form: {
        theme: '',
        cardCount: 10,
        style: 'simple',
      },
      cards: [],
      pdfUrl: '',
    };
  },
  methods: {
    async generateGame() {
      try {
        const response = await fetch('http://localhost:8080/api/generate', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(this.form),
        });
        const data = await response.json();
        this.cards = data.cards;
        this.pdfUrl = data.pdfUrl;
      } catch (error) {
        console.error('生成失敗:', error);
        alert('生成失敗，請稍後重試');
      }
    },
  },
};
</script>

<style>
@import 'tailwindcss/tailwind.css';
</style>