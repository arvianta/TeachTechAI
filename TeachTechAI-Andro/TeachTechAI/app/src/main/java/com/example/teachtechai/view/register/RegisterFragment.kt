package com.example.teachtechai.view.register

import android.os.Bundle
import android.text.Editable
import android.text.Spannable
import android.text.SpannableString
import android.text.TextWatcher
import android.text.style.ForegroundColorSpan
import android.util.Log
import android.util.Patterns
import android.util.TypedValue
import androidx.fragment.app.Fragment
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Button
import android.widget.CheckBox
import android.widget.EditText
import androidx.appcompat.app.AlertDialog
import androidx.core.content.ContextCompat
import androidx.fragment.app.viewModels
import com.example.teachtechai.R
import com.example.teachtechai.databinding.FragmentRegisterBinding
import com.example.teachtechai.view.inputotp.VerifyOTP

class RegisterFragment : Fragment() {
    private lateinit var binding : FragmentRegisterBinding
    private val registerViewModel : RegisterViewModel by viewModels()
    val verifyOTP = VerifyOTP()

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        // Inflate the layout for this fragment
        binding = FragmentRegisterBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        emailpasswordValidation()
        buttonRegister()
        observeData()
    }

    private fun observeData() {
        registerViewModel.otpResponse.observe(viewLifecycleOwner){response ->
            if(response.status == true){
                navigateToVerifyOTP(verifyOTP)
            }
        }
    }

    private fun buttonRegister() {
        binding.registerBtnDaftar.setOnClickListener {
            val name = binding.registerEditNama.text.toString()
            val email = binding.registerEditEmail.text.toString()
            val password = binding.registerEditPassword.text.toString()

            val bundle = Bundle()

            bundle.putString("email", email)
            verifyOTP.arguments = bundle
            registerViewModel.registerUser(name, email, password)
            registerViewModel.sendOTP(email)
        }
    }

    private fun emailpasswordValidation() {
        val namaEditText = binding.registerEditNama
        val emailEditText = binding.registerEditEmail
        val passwordEditText = binding.registerEditPassword
        val confirmPassEditText = binding.registerEditConfirmPassword
        val checkboxAgree = binding.registerCheckbox
        val buttonSubmit = binding.registerBtnDaftar

        val textWatcher = object : TextWatcher {
            override fun afterTextChanged(s: Editable?) {
                updateButtonState(namaEditText,emailEditText, passwordEditText, confirmPassEditText, checkboxAgree, buttonSubmit)
            }

            override fun beforeTextChanged(s: CharSequence?, start: Int, count: Int, after: Int) {}
            override fun onTextChanged(s: CharSequence?, start: Int, before: Int, count: Int) {}
        }
        emailEditText.addTextChangedListener(textWatcher)
        passwordEditText.addTextChangedListener(textWatcher)
        confirmPassEditText.addTextChangedListener(textWatcher)
        checkboxAgree.setOnCheckedChangeListener { _, _ ->
            updateButtonState(namaEditText,emailEditText, passwordEditText, confirmPassEditText, checkboxAgree, buttonSubmit)
        }

        setSpannableString(checkboxAgree)
    }
    private fun setSpannableString(checkBoxAgree: CheckBox) {
        val text = checkBoxAgree.text.toString()
        val spannableString = SpannableString(text)

        val termsAndConditionsStart = text.indexOf("Syarat dan Ketentuan")
        val termsAndConditionsEnd = termsAndConditionsStart + "Syarat dan Ketentuan".length

        val privacyPolicyStart = text.indexOf("Kebijakan Privasi")
        val privacyPolicyEnd = privacyPolicyStart + "Kebijakan Privasi".length

        val colorRed = ContextCompat.getColor(requireContext(), R.color.kaizen_primary)
        if (termsAndConditionsStart != -1 && termsAndConditionsEnd != -1) {
            spannableString.setSpan(
                ForegroundColorSpan(colorRed),
                termsAndConditionsStart,
                termsAndConditionsEnd,
                Spannable.SPAN_EXCLUSIVE_EXCLUSIVE
            )
        }

        if (privacyPolicyStart != -1 && privacyPolicyEnd != -1) {
            spannableString.setSpan(
                ForegroundColorSpan(colorRed),
                privacyPolicyStart,
                privacyPolicyEnd,
                Spannable.SPAN_EXCLUSIVE_EXCLUSIVE
            )
        }

        checkBoxAgree.text = spannableString
    }
    private fun validateEmail(emailEditText : EditText) {
        val email = emailEditText.text.toString()
        if(email.isEmpty()){
            emailEditText.error = null
        }else if (!Patterns.EMAIL_ADDRESS.matcher(email).matches()) {
            emailEditText.error = "Email tidak valid"
        } else {
            emailEditText.error = null
        }
    }
    private fun validatePassword(passwordEditText : EditText) {
        val password = passwordEditText.text.toString()
        if(password.isEmpty()){
            passwordEditText.error = null
        }else if (password.length < 8) {
            passwordEditText.error = "Kata sandi minimal 8 karakter"
        } else {
            passwordEditText.error = null
        }
    }
    private fun updateButtonState(
        editTextNama: EditText,
        editTextEmail: EditText,
        editTextPassword: EditText,
        editTextConfirmPassword: EditText,
        checkBoxAgree: CheckBox,
        buttonDaftar: Button
    ) {
        val isFormFilled = editTextNama.text.isNotEmpty() &&
                editTextEmail.text.isNotEmpty() &&
                editTextPassword.text.isNotEmpty() &&
                editTextConfirmPassword.text.isNotEmpty() &&
                checkBoxAgree.isChecked

        if (isFormFilled) {
            buttonDaftar.isEnabled = true
            buttonDaftar.background = requireContext().getDrawable(R.drawable.button_shape)
        } else {
            buttonDaftar.isEnabled = false
            buttonDaftar.background = requireContext().getDrawable(R.drawable.button_shapedisable)
        }
    }

    private fun navigateToVerifyOTP(verifyOTP: VerifyOTP){
        parentFragmentManager.beginTransaction()
            .replace(R.id.fragment_container, verifyOTP)
            .addToBackStack(null)
            .commit()
    }
}
