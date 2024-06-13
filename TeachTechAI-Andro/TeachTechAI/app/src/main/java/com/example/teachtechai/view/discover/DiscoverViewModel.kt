package com.example.teachtechai.view.discover

import android.util.Log
import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import com.example.teachtechai.data.User
import com.example.teachtechai.data.pref.UserPreference
import com.example.teachtechai.data.response.GetMeData
import com.example.teachtechai.data.response.GetMeResponse
import com.example.teachtechai.data.retrofit.ApiConfig
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response

class DiscoverViewModel : ViewModel(){
    private val _getMeResponse = MutableLiveData<GetMeResponse>()
    val getMeResponse : MutableLiveData<GetMeResponse> = _getMeResponse

    private val _errorMessage = MutableLiveData<String>()
    val errorMessage : MutableLiveData<String> = _errorMessage

    private val _isLoading = MutableLiveData<Boolean>()
    val isLoading : MutableLiveData<Boolean> = _isLoading

    private val _user = MutableLiveData<User>()
    val user : LiveData<User> = _user
    fun getMe(token : String){
        _isLoading.value = true
        val call = ApiConfig.getApiService().getme("Bearer $token")
        call.enqueue(object : Callback<GetMeResponse>{
            override fun onResponse(call: Call<GetMeResponse>, response: Response<GetMeResponse>) {
                if(response.isSuccessful){
                    _isLoading.value = false
                    _getMeResponse.value = response.body()
                }else{
                    _errorMessage.value = response.message()
                }
            }

            override fun onFailure(call: Call<GetMeResponse>, t: Throwable) {
                _errorMessage.value = t.message
            }
        })
    }
}