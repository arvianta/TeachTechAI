package com.example.teachtechai.view.editprofile

import android.app.Activity
import android.app.DatePickerDialog
import android.content.Intent
import android.net.Uri
import android.os.Bundle
import android.text.Editable
import android.util.Log
import android.util.TypedValue
import androidx.fragment.app.Fragment
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Button
import androidx.appcompat.app.AlertDialog
import androidx.fragment.app.viewModels
import androidx.navigation.fragment.findNavController
import com.bumptech.glide.Glide
import com.bumptech.glide.load.resource.bitmap.CenterCrop
import com.bumptech.glide.load.resource.bitmap.CircleCrop
import com.example.teachtechai.R
import com.example.teachtechai.data.pref.UserPreference
import com.example.teachtechai.data.pref.dataStore
import com.example.teachtechai.databinding.FragmentEditProfileBinding
import kotlinx.coroutines.runBlocking
import java.text.SimpleDateFormat
import java.util.Calendar
import java.util.Locale
import java.util.TimeZone

class EditProfileFragment : Fragment() {
    private lateinit var binding: FragmentEditProfileBinding
    private val editProfileViewModel : EditProfileViewModel by viewModels()
    private lateinit var userPreference: UserPreference
    private var selectCode = 101
    private var currentImageUri : Uri? = null
    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        // Inflate the layout for this fragment
        binding = FragmentEditProfileBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        userPreference = UserPreference.getInstance(requireContext().dataStore)
        showDatePickerDialog()
        openGallery()
        updateUser()
        uploadProfile()
        observeData()
    }
    private fun uploadProfile(){
        runBlocking {
            val token = userPreference.getToken()
            if (token != null) {
                editProfileViewModel.uploadProfile(token, currentImageUri)
            }
        }
    }
    private fun openGallery() {
        binding.editImage.setOnClickListener{
            val intent = Intent(Intent.ACTION_GET_CONTENT)
            intent.type = "image/*"
            startActivityForResult(intent, selectCode)
        }
    }

    override fun onActivityResult(requestCode: Int, resultCode: Int, data: Intent?) {
        super.onActivityResult(requestCode, resultCode, data)
        if(resultCode == Activity.RESULT_OK){
            when(requestCode){
                selectCode ->{
                    val selectedImageUri = data?.data
                    currentImageUri = selectedImageUri
                }
            }
        }
    }

    private fun showImage() {
        Glide.with(this)
            .load(currentImageUri)
            .transform(CircleCrop(), CenterCrop())
            .into(binding.editProfile)
    }

    private fun observeData() {
        editProfileViewModel.updateUserResponse.observe(viewLifecycleOwner){response->
            if(response.status == true){
                showDialogBoxUpdate()
            }
        }
        editProfileViewModel.uploadProfileResponse.observe(viewLifecycleOwner){uploadResponse->
            if(uploadResponse.status == true){
                showImage()
            }
        }
    }
    private fun updateUser(){
        runBlocking {
            val token = userPreference.getToken()
            binding.editButtonSave.setOnClickListener {
                val name = binding.editName.text.toString()
                val asal_instansi = binding.editAsalInstansi.text.toString()
                val datebirth = binding.editTanggalLahir.text.toString()
                if (token != null) {
                    editProfileViewModel.updateUser(token, name, asal_instansi, datebirth)
                }
            }
        }
    }

    private fun showDatePickerDialog() {
        val birthdate = binding.editTanggalLahir
        birthdate.setOnFocusChangeListener{_, hasFocus->
            if(hasFocus){
                val calendar = Calendar.getInstance()
                val year = calendar.get(Calendar.YEAR)
                val month = calendar.get(Calendar.MONTH)
                val day = calendar.get(Calendar.DAY_OF_MONTH)

                val datePickerDialog = DatePickerDialog(
                    requireContext(),
                    { _, selectedYear, selectedMonth, selectedDay ->
                        val selectedDate = Calendar.getInstance()
                        selectedDate.set(selectedYear, selectedMonth, selectedDay)
                        val dateFormat = SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss'Z'", Locale.getDefault())
                        dateFormat.timeZone = TimeZone.getTimeZone("UTC")
                        val formattedDate = dateFormat.format(selectedDate.time)
                        birthdate.text = Editable.Factory.getInstance().newEditable(formattedDate)
                    },
                    year, month, day
                )
                datePickerDialog.show()
            }
        }
    }

    private fun showDialogBoxUpdate(){
        val dialogView = LayoutInflater.from(requireContext()).inflate(R.layout.dialog_update_success, null)
        val dialogBuilder = AlertDialog.Builder(requireContext())
            .setView(dialogView)

        val alertDialog = dialogBuilder.create()
        alertDialog.window?.setBackgroundDrawableResource(android.R.color.transparent)
        alertDialog.show()

        val width = TypedValue.applyDimension(TypedValue.COMPLEX_UNIT_DIP, 301f, resources.displayMetrics).toInt()
        val height = TypedValue.applyDimension(TypedValue.COMPLEX_UNIT_DIP, 250f, resources.displayMetrics).toInt()
        alertDialog.window?.setLayout(width, height)
        val buttonSuccess = dialogView.findViewById<Button>(R.id.buttonOk)
        buttonSuccess.setOnClickListener {
            alertDialog.dismiss()
            navigateToProfile()
        }
    }

    private fun navigateToProfile(){
        findNavController().navigate(R.id.editProfile_to_profileFragment)
    }
}