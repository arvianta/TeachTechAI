package com.example.teachtechai.view.chat

import android.util.Log
import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import com.example.teachtechai.data.pref.UserPreference
import com.example.teachtechai.data.response.ChatResponse
import com.example.teachtechai.data.retrofit.ApiConfig
import com.example.teachtechai.data.retrofit.ApiService
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response

class ChatViewModel : ViewModel() {
    private val _chatMessage = MutableLiveData<ChatResponse>()
    val chatMessage : LiveData<ChatResponse> = _chatMessage

    private val _errorMessage = MutableLiveData<String>()
    val errorMessage : LiveData<String> = _errorMessage
    fun getChatResponse(token : String, prompt: String) {
        val call = ApiConfig.getApiService().getChatResponse("Bearer $token", "Menentukan Metode Pembelajaran Di Kelas", prompt, "TeachTechAI")
        call.enqueue(object : Callback<ChatResponse> {
            override fun onResponse(call: Call<ChatResponse>, response: Response<ChatResponse>) {
                if(response.isSuccessful){
                    Log.d("TES", "${response.body()}")
                    _chatMessage.value = response.body()
                }else{
                    _errorMessage.value = response.message()
                }
            }

            override fun onFailure(call: Call<ChatResponse>, t: Throwable) {
                _errorMessage.value = t.message
            }
        })
    }
}