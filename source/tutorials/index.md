---
title: 教程
layout: docs
---

<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.1/css/all.min.css">

<div class="tutorial-header">
  <h1><i class="fa-solid fa-graduation-cap"></i> 教程合集</h1>
  <p class="tutorial-subtitle">记录学习路上的每一步，从入门到实践</p>
</div>

<div class="tutorial-grid">

  <a href="/typescript-tutorial/" class="tutorial-card">
    <div class="card-icon">
      <i class="fa-brands fa-typescript"></i>
    </div>
    <div class="card-body">
      <h2>TypeScript 快速入门</h2>
      <p>从零开始的 TypeScript 教程，涵盖基础类型、接口、泛型、高级类型等核心概念。内含在线代码编辑器，可直接运行示例代码。</p>
      <div class="card-tags">
        <span class="tag">TypeScript</span>
        <span class="tag">入门</span>
        <span class="tag">前端</span>
      </div>
      <span class="card-link">开始学习 <i class="fa-solid fa-arrow-right"></i></span>
    </div>
  </a>

  <a href="/python-tutorial/" class="tutorial-card">
    <div class="card-icon">
      <i class="fa-brands fa-python"></i>
    </div>
    <div class="card-body">
      <h2>Python 快速入门</h2>
      <p>从零开始的 Python 教程，涵盖基础数据类型、控制流、函数、面向对象、列表推导式等核心概念。内含在线代码编辑器，可直接运行 Python 代码。</p>
      <div class="card-tags">
        <span class="tag">Python</span>
        <span class="tag">入门</span>
        <span class="tag">后端</span>
      </div>
      <span class="card-link">开始学习 <i class="fa-solid fa-arrow-right"></i></span>
    </div>
  </a>

</div>

<style>
.tutorial-header {
  text-align: center;
  padding: 40px 0 30px;
}
.tutorial-header h1 {
  font-size: 2em;
  margin-bottom: 10px;
  color: var(--color-h1);
}
.tutorial-header h1 i {
  margin-right: 10px;
  color: var(--theme-color, #8b3a3a);
}
.tutorial-subtitle {
  color: var(--color-meta);
  font-size: 1.05em;
}

.tutorial-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  gap: 24px;
  padding: 10px 0 40px;
}

.tutorial-card {
  display: flex;
  flex-direction: column;
  background: var(--color-card);
  border: 1px solid var(--color-block);
  border-radius: 8px;
  overflow: hidden;
  text-decoration: none;
  color: inherit;
  transition: transform 0.2s, box-shadow 0.2s;
}
.tutorial-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0,0,0,0.12);
}

.card-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 80px;
  background: linear-gradient(135deg, var(--theme-color, #8b3a3a), #a05040);
  color: #fff;
  font-size: 2.8em;
}

.card-body {
  padding: 20px;
  display: flex;
  flex-direction: column;
  flex: 1;
}
.card-body h2 {
  font-size: 1.2em;
  margin: 0 0 10px;
  color: var(--color-h2);
}
.card-body p {
  font-size: 0.93em;
  line-height: 1.6;
  color: var(--color-p);
  margin: 0 0 14px;
  flex: 1;
}

.card-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 14px;
}
.tag {
  display: inline-block;
  font-size: 0.78em;
  padding: 3px 10px;
  border-radius: 4px;
  background: var(--color-block);
  color: var(--color-meta);
}

.card-link {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 0.9em;
  font-weight: 600;
  color: var(--theme-color, #8b3a3a);
  transition: gap 0.2s;
}
.tutorial-card:hover .card-link {
  gap: 10px;
}
.card-link i {
  font-size: 0.85em;
}

@media (max-width: 640px) {
  .tutorial-grid {
    grid-template-columns: 1fr;
  }
  .tutorial-header {
    padding: 20px 0;
  }
  .tutorial-header h1 {
    font-size: 1.5em;
  }
}
</style>
