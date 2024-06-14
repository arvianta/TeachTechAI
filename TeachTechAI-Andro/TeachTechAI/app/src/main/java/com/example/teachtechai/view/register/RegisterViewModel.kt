package com.example.teachtechai.view.register

import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.example.teachtechai.data.response.ApiError
import com.example.teachtechai.data.response.OTPResponse
import com.example.teachtechai.data.response.RegisterResponse
import com.example.teachtechai.data.retrofit.ApiConfig
import com.google.gson.Gson
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch
import retrofit2.Call
import retrofit2.Response
import javax.security.auth.callback.Callback

class RegisterViewModel : ViewModel() {
    private val _registerResult = MutableLiveData<RegisterResponse>()
    val registerResult : LiveData<RegisterResponse> = _registerResult

    private val _isLoading = MutableLiveData<Boolean>()
    val isLoading : LiveData<Boolean> = _isLoading

    private val _errorMessage = MutableLiveData<String?>()
    val errorMessage : LiveData<String?> = _errorMessage

    private val _errorStatus = MutableLiveData<Boolean>()
    val errorStatus : LiveData<Boolean> = _errorStatus

    private val _otpResponse = MutableLiveData<OTPResponse>()
    val otpResponse : LiveData<OTPResponse> = _otpResponse

    private val _otpError = MutableLiveData<String?>()
    val otpError : LiveData<OTPResponse> = _otpResponse

    fun registerUser(email : String, name : String, password : String){
        val call = ApiConfig.getApiService().register(email, name, password)
        call.enqueue(object : retrofit2.Callback<RegisterResponse> {
            override fun onResponse(call: Call<RegisterResponse>, response: Response<RegisterResponse>) {
                if(response.isSuccessful){
                    _registerResult.value = response.body()
                }else{
                    response.errorBody()?.let{errorBody->
                        val errorResponse = errorBody.string()
                        val apiError = Gson().fromJson(errorResponse, ApiError::class.java)
                        _errorMessage.value = apiError.error
                    }
                }
            }

            override fun onFailure(call: Call<RegisterResponse>, t: Throwable) {
                _errorMessage.value = t.message
            }
        })
    }
}