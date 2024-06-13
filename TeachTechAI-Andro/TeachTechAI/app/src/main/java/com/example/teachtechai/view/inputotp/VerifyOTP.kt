package com.example.teachtechai.view.inputotp

import android.os.Bundle
import android.util.Log
import android.util.TypedValue
import androidx.fragment.app.Fragment
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Button
import androidx.appcompat.app.AlertDialog
import androidx.fragment.app.viewModels
import com.example.teachtechai.R
import com.example.teachtechai.databinding.FragmentVerifyOtpBinding
import com.example.teachtechai.view.login.LoginFragment

class VerifyOTP : Fragment() {
    private lateinit var binding: FragmentVerifyOtpBinding
    private val verifyViewModel : VerifyViewModel by viewModels()
    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        binding = FragmentVerifyOtpBinding.inflate(inflater, container,false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        verifyOTP()
        observeData()
    }

    private fun observeData() {
        verifyViewModel.verifyOTPResponse.observe(viewLifecycleOwner){response->
            if(response.status == true){
                showDialogBoxRegister()
            }
        }
    }

    private fun verifyOTP() {
        val data = arguments
        val email = data?.getString("email")

        binding.fpButtonVerifikasi.setOnClickListener {
            val otp = binding.verifyInputOTP.text.toString()
            if(email != null){
                verifyViewModel.verifyOTP(email, otp)
            }
        }
    }
    private fun showDialogBoxRegister(){
        val dialogView = LayoutInflater.from(requireContext()).inflate(R.layout.dialog_registration_success, null)
        val dialogBuilder = AlertDialog.Builder(requireContext())
            .setView(dialogView)

        val alertDialog = dialogBuilder.create()
        alertDialog.window?.setBackgroundDrawableResource(android.R.color.transparent)
        alertDialog.show()

        val width = TypedValue.applyDimension(TypedValue.COMPLEX_UNIT_DIP, 301f, resources.displayMetrics).toInt()
        val height = TypedValue.applyDimension(TypedValue.COMPLEX_UNIT_DIP, 315f, resources.displayMetrics).toInt()
        alertDialog.window?.setLayout(width, height)
        val buttonSuccess = dialogView.findViewById<Button>(R.id.buttonOk)
        buttonSuccess.setOnClickListener {
            alertDialog.dismiss()
            navigateToLogin()
        }
    }

    private fun navigateToLogin(){
        parentFragmentManager.beginTransaction()
            .replace(R.id.fragment_container, LoginFragment())
            .addToBackStack(null)
            .commit()
    }
}