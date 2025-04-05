<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <h1>服装进销存系统</h1>
        <p>HD-PSI</p>
      </div>
      
      <n-form
        ref="formRef"
        :model="formValue"
        :rules="rules"
        label-placement="left"
        size="large"
        class="login-form"
      >
        <n-form-item path="username" label="用户名">
          <n-input v-model:value="formValue.username" placeholder="请输入用户名" />
        </n-form-item>
        
        <n-form-item path="password" label="密码">
          <n-input
            v-model:value="formValue.password"
            type="password"
            show-password-on="click"
            placeholder="请输入密码"
          />
        </n-form-item>
        
        <div class="login-actions">
          <n-button
            type="primary"
            size="large"
            block
            :loading="loading"
            @click="handleLogin"
          >
            登录
          </n-button>
        </div>
      </n-form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import { NForm, NFormItem, NInput, NButton } from 'naive-ui'
import auth from '../services/auth'

const router = useRouter()
const message = useMessage()
const formRef = ref(null)
const loading = ref(false)

// 表单数据
const formValue = reactive({
  username: '',
  password: ''
})

// 表单验证规则
const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6个字符', trigger: 'blur' }
  ]
}

// 登录处理
const handleLogin = (e) => {
  e.preventDefault()
  formRef.value?.validate(async (errors) => {
    if (!errors) {
      loading.value = true
      try {
        const response = await auth.login(formValue.username, formValue.password)
        
        // 保存令牌和用户信息
        localStorage.setItem('token', response.token)
        localStorage.setItem('user', JSON.stringify(response.user))
        
        message.success('登录成功')
        router.push('/dashboard')
      } catch (error) {
        console.error('登录失败:', error)
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f0f2f5;
}

.login-card {
  width: 400px;
  padding: 40px;
  background-color: #fff;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-header h1 {
  font-size: 24px;
  color: #18a058;
  margin-bottom: 8px;
}

.login-header p {
  font-size: 14px;
  color: #999;
}

.login-form {
  margin-top: 20px;
}

.login-actions {
  margin-top: 30px;
}
</style>
