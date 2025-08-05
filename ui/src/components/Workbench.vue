<script>
import { ref, onMounted } from 'vue'
import Separator from "./ui/separator/Separator.vue";
import { apiFetch } from "@/api/request";

export default {
  components: {
    Separator
  },

  setup() {
    const userInput = ref('')
    const response = ref('')
    const error = ref('')
    const loading = ref(false)
    const messages = ref([]);
    const schema = ref('');

    // Fetch the schema
    apiFetch("/cdb-schema").then(async res => {
      if (res.ok) {
        schema.value = await res.text();

        messages.value.push({"role": "system", "content": `
Your task is to convert a question into a SQL query, given a Postgres database schema.
Adhere to these rules:
- **Deliberately go through the question and database schema word by word** to appropriately answer the question
- **Use Table Aliases** to prevent ambiguity. For example, \`SELECT table1.col1, table2.col1 FROM table1 JOIN table2 ON table1.id = table2.id\`.
- When creating a ratio, always cast the numerator as float
- The assistant's response MUST be programmatic. The response must be pure SQL, or HTML for a table if the user requests a list of tables, parameters, etc. The response MUST NOT contain any prefix or suffix, and must be pure SQL or HTML.
- The request MUST not contain any additional text, comments, or explanations. The response MUST be a valid SQL query that can be executed on the database schema provided above. The assistant MUST NOT add any prefix or suffix into their responses, such as '###', 'Response:', or anything else. Only the SQL wrapped in \`\`\` tags is acceptable.
- When thinking, thoughts are placed in the XML tag <thought>...</thought>. The final result as an SQL string is placed in <final>...</final>

\`\`\`ASSISTANT INSTRUCTIONS FOR THE FUTURE - DO NOT INCLUDE THIS IN THE FINAL RESPONSE\`\`\`
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
        messages.value.push({"role":"user","content":userInput.value})
      }
      
      userInput.value = '';

      // Send the message history to ollama
      apiFetch(API_URL, 'POST', JSON.stringify({"messages":messages.value})).then(async (res) => {
        // Mark that we are finished loading.
        loading.value = false

        if (!res.ok) {
          // In the case of an error, push the error as a message
          let err = `Error ${res.status}: ${res.statusText}`;
          messages.value.push({"role": "error", "content": err});
          throw new Error(err);

        } else {
          // In the case of a success, push the response
          const data = await res.json()

          // Append the message to the message list
          messages.value.push(data)
        }
      });
    }

    return {
      userInput,
      response,
      error,
      loading,
      sendQuery,
      messages
    }
  }
}
</script>

<template>
  <div class="h-full flex flex-col">
    <div class="p-2 flex justify-between bg-accent">
      <p>AI Chat</p>
    </div>

    <Separator></Separator>
    
    <div class="p-2 flex-2 max-h-[calc(100vh-180px)] overflow-y-scroll overflow-x-hidden">
      <div v-for="(message, index) in messages" :key="index" class="text-sm flex flex-col items-start">
        <p v-if="message.role != 'system'" class="block p-1 font-black" :class="message.role == 'user' ? 'self-end' : 'self-start'">
          {{message.role == 'user' ? "You:" : "StrataBot:"}}
        </p>

        <div v-if="message.role != 'system'"  class="w-fit rounded-sm p-2 max-w-4/5"
        :class="message.role == 'user' ? 'bg-blue-300 self-end' : message.role == 'assistant' ? 'bg-purple-300 self-start' : 'bg-red-400'">
          <p class="block">{{message.content}}</p>
        </div>
      </div>
    </div>

    <Separator></Separator>

    <div v-if="loading" class="flex justify-items-center">
      <p class="block">Thinking...</p>
    </div>

    <div class="py-3 px-4 flex items-center input-box">
      <textarea placeholder="Write a message..." v-model="userInput"
        class="w-full resize-none break-words focus:outline-none" rows="1" ref="chatBox"
        @input="scaleChatBox">
      </textarea>
      
      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
        stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
        class="lucide lucide-send-horizontal-icon lucide-send-horizontal"
        :class="userInput == '' ? 'text-neutral-400' : 'text-black cursor-pointer'"
        @click="sendQuery" :disabled="loading">
        <path
          d="M3.714 3.048a.498.498 0 0 0-.683.627l2.843 7.627a2 2 0 0 1 0 1.396l-2.842 7.627a.498.498 0 0 0 .682.627l18-8.5a.5.5 0 0 0 0-.904z" />
        <path d="M6 12h16" />
      </svg>
    </div>
  </div>
</template>
