package com.example.teachtechai.view.changepassword

import android.os.Bundle
import android.util.TypedValue
import androidx.fragment.app.Fragment
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Button
import androidx.appcompat.app.AlertDialog
import androidx.fragment.app.viewModels
import androidx.navigation.fragment.findNavController
import com.example.teachtechai.R
import com.example.teachtechai.data.pref.UserPreference
import com.example.teachtechai.data.pref.dataStore
import com.example.teachtechai.databinding.FragmentChangePasswordBinding
import kotlinx.coroutines.runBlocking

class ChangePasswordFragment : Fragment() {
    private lateinit var binding: FragmentChangePasswordBinding
    private lateinit var userPreference : UserPreference
    private val changePasswordViewModel : ChangePasswordViewModel by viewModels()
    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        // Inflate the layout for this fragment
        binding = FragmentChangePasswordBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        userPreference = UserPreference.getInstance(requireContext().dataStore)
        changePassword()
        observeData()
    }

    private fun observeData() {
        changePasswordViewModel.changeResponse.observe(viewLifecycleOwner){response->
            if(response.status == true){
                showDialogBoxChange()
            }
        }
    }

    private fun changePassword() {
        runBlocking {
            val token = userPreference.getToken()
            binding.editButtonSave.setOnClickListener {
                val old_password = binding.cpOldPassword.text.toString()
                val new_password = binding.cpNewPassword.text.toString()

                if(token != null){
                    changePasswordViewModel.changePassword(token, old_password, new_password)
                }
            }
        }
    }

    private fun showDialogBoxChange(){
        val dialogView = LayoutInflater.from(requireContext()).inflate(R.layout.dialog_changepassword_success, null)
        val dialogBuilder = AlertDialog.Builder(requireContext())
            .setView(dialogView)

        val alertDialog = dialogBuilder.create()
        alertDialog.window?.setBackgroundDrawableResource(android.R.color.transparent)
        alertDialog.show()

        val width = TypedValue.applyDimension(TypedValue.COMPLEX_UNIT_DIP, 301f, resources.displayMetrics).toInt()
        val height = TypedValue.applyDimension(TypedValue.COMPLEX_UNIT_DIP, 301f, resources.displayMetrics).toInt()
        alertDialog.window?.setLayout(width, height)
        val buttonSuccess = dialogView.findViewById<Button>(R.id.buttonOk)
        buttonSuccess.setOnClickListener {
            alertDialog.dismiss()
            navigateToProfile()
        }
    }

    private fun navigateToProfile(){
        findNavController().navigate(R.id.changePasswordFragment_to_profileFragment)
    }
}