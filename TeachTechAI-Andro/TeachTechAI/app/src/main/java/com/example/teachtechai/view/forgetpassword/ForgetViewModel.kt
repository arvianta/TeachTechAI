package com.example.teachtechai.view.forgetpassword

import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import com.example.teachtechai.data.response.ChangePasswordResponse
import com.example.teachtechai.data.retrofit.ApiConfig
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response

class ForgetViewModel : ViewModel() {
    private val _forgetPasswordResponse = MutableLiveData<ChangePasswordResponse>()
    val forgetPasswordResponse : LiveData<ChangePasswordResponse> = _forgetPasswordResponse

    private val _errorMessage = MutableLiveData<String>()
    val errorMessage : LiveData<String> = _errorMessage

    fun forgetPassword(email : String){
        val call = ApiConfig.getApiService().forgetpassword(email)
        call.enqueue(object : Callback<ChangePasswordResponse> {
            override fun onResponse(
                call: Call<ChangePasswordResponse>,
                response: Response<ChangePasswordResponse>
            ) {
                if(response.isSuccessful){
                    _forgetPasswordResponse.value = response.body()
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