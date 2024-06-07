package com.example.teachtechai.view.login

import android.util.Log
import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.example.teachtechai.data.UserRepository
import com.example.teachtechai.data.pref.UserModel
import com.example.teachtechai.data.response.LoginResponse
import com.example.teachtechai.data.retrofit.ApiConfig
import kotlinx.coroutines.launch
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response


class LoginViewModel (private val repository: UserRepository) : ViewModel() {
    private val _loginResult = MutableLiveData<LoginResponse>()
    val loginResult : LiveData<LoginResponse> = _loginResult

    private val _isLoading = MutableLiveData<Boolean>()
    val isLoading : LiveData<Boolean> = _isLoading

    private val _errorMessage = MutableLiveData<String?>()
    val errorMessage : LiveData<String?> = _errorMessage

    fun loginUser(email: String, password: String) {
        _isLoading.value = true
        val call = ApiConfig.getApiService().login(email, password)
        call.enqueue(object : Callback<LoginResponse> {
            override fun onResponse(call: Call<LoginResponse>, response: Response<LoginResponse>) {
                if (response.isSuccessful) {
                    _loginResult.value = response.body()
                    _isLoading.value = false
                    val token = response.body()?.data?.sessionToken
                    if(token != null){
                        saveSession(UserModel(email, token, true))
                    }else{
                        _errorMessage.value = "Token is null"
                    }
                }else{
                    _isLoading.value = false
                    _errorMessage.value = response.body()?.message.toString()
                }
            }

            override fun onFailure(call: Call<LoginResponse>, t: Throwable) {
                _errorMessage.value = t.message?: "Login User Error"
            }
        })
    }
    private fun saveSession(user: UserModel){
        viewModelScope.launch{
            repository.saveSession(user)
            Log.d("Login View Model", "Save Session : {$user}" )
        }
    }
}