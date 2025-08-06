<script setup>
import { ref, onMounted, nextTick } from 'vue'
import Separator from "./ui/separator/Separator.vue";
import { apiFetch } from "@/api/request";
import hljs from 'highlight.js';
import MarkdownIt from 'markdown-it';

const userInput = ref('')
const response = ref('')
const error = ref('')
const loading = ref(false)
const messages = ref([]);
const schema = ref('');

const md = MarkdownIt({
  html: true,
  linkify: true,
  typographer: true
})
// Fetch the schema
apiFetch("/cdb-schema").then(async res => {
  if (res.ok) {
    schema.value = await res.text();

    messages.value.push({
      "role": "system", "content": `
      You are a professional database engineer tasked with helping the user to read and understand data from their postgres database.
Your task is to produce SQL queries and/or explanations to satisfy their request using the provided Postgres database schema.
Adhere to these rules:
- **Deliberately go through the question and database schema word by word** to appropriately answer the question
- **Use Table Aliases** to prevent ambiguity. For example, \`SELECT table1.col1, table2.col1 FROM table1 JOIN table2 ON table1.id = table2.id\`.
- Your explanations must be simple and consise, as your intended audience may be people with little to no technical skill.  

\`\`\`ASSISTANT INSTRUCTIONS FOR THE FUTURE - DO NOT INCLUDE THIS IN THE FINAL RESPONSE\`\`\`
If the user does not start with a question or actionable message, respond with a short friendly message offering to help explain something or build a query
Respond to the user in a friendly way, but minimize any extra words. Answer their inquirty directly and consisely.
Generate a SQL query that answers the question \`{question}\`.
This query will run on a database whose schema is represented by this string:
\`\`\`sql
${schema.value}
\`\`\`sql
`})
  } else {
    throw new Error(`Error ${res.status}: ${res.statusText}`);
  }
});

const API_URL = '/stratabot'
const sendQuery = async () => {
  if (!userInput.value.trim()) return

  loading.value = true
  response.value = ''
  error.value = ''

  // if the messages list is empty or the last message was not sent by the user...
  if (messages.value.length == 0 || messages.value[messages.value.length - 1].role != "user") {
    messages.value.push({ "role": "user", "content": userInput.value })
  }

  userInput.value = '';
  scaleChatBox();

  // Send the message history to ollama
  apiFetch(API_URL, 'POST', JSON.stringify({ "messages": messages.value })).then(async (res) => {
    // Mark that we are finished loading.
    loading.value = false

    if (!res.ok) {
      // In the case of an error, push the error as a message
      let err = `Error ${res.status}: ${res.statusText}`;
      messages.value.push({ "role": "error", "content": err });
      throw new Error(err);

    } else {
      // In the case of a success, push the response
      let data = await res.json()
      data.pretty = md.render(data.content)

      // Append the message to the message list
      messages.value.push(data)
    }
  });


}

const inputBox = ref(null);

function scaleChatBox() {
  nextTick(() => {
    const chatBox = inputBox.value;
    if (chatBox) {
      chatBox.style.height = 'auto'; // Reset height
      chatBox.style.height = `${chatBox.scrollHeight}px`; // Set to scrollHeight
    }
  });
}

function handleEnter(e) {
  if (e.shiftKey) {
    // Allow newline
    userInput.value += '\n';
    scaleChatBox();
    return
  }

  // Send message on Enter
  e.preventDefault()
  sendQuery()
}
</script>

<template>
  <div class="h-full flex flex-col">
    <div class="p-2 flex justify-between bg-accent">
      <p>AI Chat</p>
    </div>

    <Separator></Separator>

    <div class="p-2 flex-1 overflow-y-scroll overflow-x-hidden">
      <div v-for="(message, index) in messages" :key="index" class="flex flex-col items-start mx-2 *:my-3">
        <p v-if="message.role == 'user'" class="self-end p-2 px-3 bg-accent rounded-md whitespace-pre-wrap">{{
          message.content }}</p>
        <p v-if="message.role == 'assistant'" v-html="message.pretty"
          class="prose prose-pre:bg-accent prose-pre:border-1 prose-pre:text-black"></p>


        <!-- <p v-if="message.role != 'system'" class="block p-1 font-black"
          :class="message.role == 'user' ? 'self-end' : 'self-start'">
          {{ message.role == 'user' ? "You:" : "StrataBot:" }}
        </p>

        <div v-if="message.role != 'system'" class="w-fit rounded-sm p-2 max-w-4/5"
          :class="message.role == 'user' ? 'bg-blue-300 self-end' : message.role == 'assistant' ? 'bg-purple-300 self-start' : 'bg-red-400'">
          <p class="block">{{ message.content }}</p>
        </div> -->
      </div>
    </div>

    <Separator></Separator>

    <div v-if="loading" class="flex justify-items-center">
      <p class="block">Thinking...</p>
    </div>

    <div class="py-3 px-4 flex items-end gap-2">
      <textarea placeholder="Write a message..." v-model="userInput" @keydown.enter.prevent="handleEnter"
        class="w-full resize-none break-words focus:outline-none" rows="1" ref="inputBox" @input="scaleChatBox()">
      </textarea>

      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
        stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
        class="lucide lucide-send-horizontal-icon lucide-send-horizontal"
        :class="userInput == '' ? 'text-neutral-400' : 'text-black cursor-pointer'" @click="sendQuery"
        :disabled="loading">
        <path
          d="M3.714 3.048a.498.498 0 0 0-.683.627l2.843 7.627a2 2 0 0 1 0 1.396l-2.842 7.627a.498.498 0 0 0 .682.627l18-8.5a.5.5 0 0 0 0-.904z" />
        <path d="M6 12h16" />
      </svg>
    </div>
  </div>
</template>
