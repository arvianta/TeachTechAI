package com.example.teachtechai.view.chat

import android.adservices.topics.Topic
import android.util.Log
import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import com.example.teachtechai.data.response.TopicItem
import com.example.teachtechai.data.response.TopicResponse
import com.example.teachtechai.data.retrofit.ApiConfig
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response

class DrawerViewModel : ViewModel() {
    private val _topicResponse = MutableLiveData<List<TopicItem>?>()
    val topicResponse : LiveData<List<TopicItem>?> = _topicResponse

    private val _errorMessage = MutableLiveData<String>()
    val errorMessage : LiveData<String> = _errorMessage

    fun getTopic(token : String){
        val call = ApiConfig.getApiService().getAllConversations("Bearer $token")
        call.enqueue(object : Callback<TopicResponse>{
            override fun onResponse(call: Call<TopicResponse>, response: Response<TopicResponse>) {
                if(response.isSuccessful){
                    _topicResponse.value = response.body()?.data?.filterNotNull()
                }else{
                    _errorMessage.value = response.message()
                }
            }

            override fun onFailure(call: Call<TopicResponse>, t: Throwable) {
                _errorMessage.value = t.message
            }

        })
    }
}