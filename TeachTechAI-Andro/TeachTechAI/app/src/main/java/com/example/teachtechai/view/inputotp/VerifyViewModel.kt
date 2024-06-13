package com.example.teachtechai.view.inputotp

import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import com.example.teachtechai.data.response.VerifyOTPResponse
import com.example.teachtechai.data.retrofit.ApiConfig
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response

class VerifyViewModel : ViewModel(){
    private val _verifyOTPResponse = MutableLiveData<VerifyOTPResponse>()
    val verifyOTPResponse : LiveData<VerifyOTPResponse> = _verifyOTPResponse

    private val _verifyOTPError = MutableLiveData<String>()
    val verifyOTPError : LiveData<String> = _verifyOTPError

    fun verifyOTP(email : String, otp : String){
        val call = ApiConfig.getApiService().verifyotp(email, otp)
        call.enqueue(object : Callback<VerifyOTPResponse> {
            override fun onResponse(call: Call<VerifyOTPResponse>, response: Response<VerifyOTPResponse>) {
                if(response.isSuccessful){
                    _verifyOTPResponse.value = response.body()
                }else{
                    _verifyOTPError.value = response.message()
                }
            }

            override fun onFailure(call: Call<VerifyOTPResponse>, t: Throwable) {
                _verifyOTPError.value = t.message
            }

        })
    }
}
