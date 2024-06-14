package com.example.teachtechai.view.changepassword

import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import com.example.teachtechai.data.response.ChangePasswordResponse
import com.example.teachtechai.data.retrofit.ApiConfig
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response

class ChangePasswordViewModel : ViewModel() {
    private val _changeResponse = MutableLiveData<ChangePasswordResponse>()
    val changeResponse : LiveData<ChangePasswordResponse> = _changeResponse

    private val _errorMessage = MutableLiveData<String>()
    val errorMessage : LiveData<String> = _errorMessage

    fun changePassword(token : String, old_password : String, new_password : String){
        val call = ApiConfig.getApiService().changepassword("Bearer $token", old_password, new_password)
        call.enqueue(object : Callback<ChangePasswordResponse> {
            override fun onResponse(
                call: Call<ChangePasswordResponse>,
                response: Response<ChangePasswordResponse>
            ) {
                if(response.isSuccessful){
                    _changeResponse.value = response.body()
                }else{
                    _errorMessage.value = response.message()
                }
            }

            override fun onFailure(call: Call<ChangePasswordResponse>, t: Throwable) {
                _errorMessage.value = t.message
            }
        })
    }
}