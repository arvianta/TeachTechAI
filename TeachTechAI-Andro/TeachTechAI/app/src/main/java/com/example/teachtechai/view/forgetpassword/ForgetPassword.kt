package com.example.teachtechai.view.forgetpassword

import android.os.Bundle
import android.util.TypedValue
import androidx.fragment.app.Fragment
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Button
import androidx.appcompat.app.AlertDialog
import androidx.fragment.app.viewModels
import com.example.teachtechai.R
import com.example.teachtechai.databinding.FragmentForgetPasswordBinding
import com.example.teachtechai.view.inputotp.VerifyOTP
import com.example.teachtechai.view.login.LoginFragment

class ForgetPassword : Fragment() {
    private lateinit var binding : FragmentForgetPasswordBinding
    private val fpViewModel : ForgetViewModel by viewModels()
    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        // Inflate the layout for this fragment
        binding = FragmentForgetPasswordBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        resetPassword()
        observeData()
    }

    private fun observeData() {
        fpViewModel.forgetPasswordResponse.observe(viewLifecycleOwner){response ->
            if(response.status == true){
                showDialogBoxReset()
            }
        }
    }

    private fun resetPassword() {
        binding.fpButtonReset.setOnClickListener {
            val email = binding.fpEditEmail.text.toString()
            fpViewModel.forgetPassword(email)
        }
    }

    private fun navigateToLogin(){
        parentFragmentManager.beginTransaction()
            .replace(R.id.fragment_container, LoginFragment())
            .addToBackStack(null)
            .commit()
    }

    private fun showDialogBoxReset(){
        val dialogView = LayoutInflater.from(requireContext()).inflate(R.layout.dialog_forgetpassword_success, null)
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
}