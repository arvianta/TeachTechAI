package com.example.teachtechai.view.profile

import android.util.Log
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.example.teachtechai.data.UserRepository
import com.example.teachtechai.data.response.LogoutResponse
import com.example.teachtechai.data.retrofit.ApiConfig
import kotlinx.coroutines.launch
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response

class LogoutViewModel(private val repository: UserRepository) : ViewModel() {
    fun logoutUser(token : String){
        val call = ApiConfig.getApiService().logout("Bearer $token")
        call.enqueue(object : Callback<LogoutResponse> {
            override fun onResponse(call: Call<LogoutResponse>, response: Response<LogoutResponse>) {
                logoutSession()
                if (response.isSuccessful) {
                    val logoutResponse = response.body()
                    if (logoutResponse != null) {
                        // Handle successful logout
                        println("Logout successful: ${logoutResponse.message}")
                    } else {
                        // Handle logout failure
                        println("Logout failed: ${logoutResponse?.message}")
                    }
                } else {
                    // Handle request failure
                    println("Request failed: ${response.errorBody()?.string()}")
                }
            }

            override fun onFailure(call: Call<LogoutResponse>, t: Throwable) {
                // Handle network failure
                println("Network error: ${t.message}")
            }
        })
    }
    private fun logoutSession(){
        viewModelScope.launch{
            repository.logout()
        }
    }
}