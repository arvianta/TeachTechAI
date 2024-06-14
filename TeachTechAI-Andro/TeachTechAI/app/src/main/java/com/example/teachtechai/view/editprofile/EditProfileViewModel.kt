package com.example.teachtechai.view.editprofile

import android.app.Application
import android.net.Uri
import androidx.lifecycle.AndroidViewModel
import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import com.example.teachtechai.data.response.ChangePasswordResponse
import com.example.teachtechai.data.response.UploadProfileResponse
import com.example.teachtechai.data.retrofit.ApiConfig
import com.example.teachtechai.view.utils.reduceFileImage
import com.example.teachtechai.view.utils.uriToFile
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.MultipartBody
import okhttp3.RequestBody.Companion.asRequestBody
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response

class EditProfileViewModel (application: Application) :AndroidViewModel(application){
    private val _updateUserResponse = MutableLiveData<ChangePasswordResponse>()
    val updateUserResponse : LiveData<ChangePasswordResponse> = _updateUserResponse

    private val _errorMessage = MutableLiveData<String>()
    val errorMessage : LiveData<String> = _errorMessage

    private val _uploadProfileResponse = MutableLiveData<UploadProfileResponse>()
    val uploadProfileResponse : LiveData<UploadProfileResponse> = _uploadProfileResponse

    private val _uploadErrorMessage = MutableLiveData<String>()
    val uploadErrorMessage : LiveData<String> = _uploadErrorMessage

    fun updateUser(token : String, name : String, asal_instansi : String, date_of_birth : String){
        val call = ApiConfig.getApiService().updateuser("Bearer $token",name, asal_instansi, date_of_birth)
        call.enqueue(object : Callback<ChangePasswordResponse> {
            override fun onResponse(
                call: Call<ChangePasswordResponse>,
                response: Response<ChangePasswordResponse>
            ) {
                if(response.isSuccessful){
                    _updateUserResponse.value = response.body()
                }else{
                    _errorMessage.value = response.message()
                }
            }

            override fun onFailure(call: Call<ChangePasswordResponse>, t: Throwable) {
                _errorMessage.value = t.message
            }
        })
    }

    fun uploadProfile(token: String, uri : Uri?){
        uri?.let{imageUri->
            val imageFile = uriToFile(imageUri, getApplication()).reduceFileImage()

            val requestImageFile = imageFile.asRequestBody("image/jpeg".toMediaType())
            val multipartBody = MultipartBody.Part.createFormData(
                "file",
                imageFile.name,
                requestImageFile
            )

            val call = ApiConfig.getApiService().uploadProfile("Bearer $token", multipartBody)
            call.enqueue(object : Callback<UploadProfileResponse>{
                override fun onResponse(
                    call: Call<UploadProfileResponse>,
                    response: Response<UploadProfileResponse>
                ) {
                    if(response.isSuccessful){
                        _uploadProfileResponse.value = response.body()
                    }else{
                        _uploadErrorMessage.value = response.message()
                    }
                }

                override fun onFailure(call: Call<UploadProfileResponse>, t: Throwable) {
                    _uploadErrorMessage.value = t.message
                }

            })
        }
    }

}