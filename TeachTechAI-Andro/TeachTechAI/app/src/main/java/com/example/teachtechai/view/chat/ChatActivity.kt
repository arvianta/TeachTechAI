package com.example.teachtechai.view.chat

import android.content.Intent
import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.widget.LinearLayout
import androidx.activity.viewModels
import androidx.core.view.GravityCompat
import androidx.drawerlayout.widget.DrawerLayout
import androidx.lifecycle.Observer
import androidx.recyclerview.widget.DividerItemDecoration
import androidx.recyclerview.widget.LinearLayoutManager
import androidx.recyclerview.widget.RecyclerView
import com.example.teachtechai.R
import com.example.teachtechai.data.Message
import com.example.teachtechai.data.pref.UserPreference
import com.example.teachtechai.data.pref.dataStore
import com.example.teachtechai.databinding.ActivityChatBinding
import com.example.teachtechai.databinding.LayoutDrawerBinding
import com.example.teachtechai.view.discover.TopicAdapter
import com.example.teachtechai.view.topic.DrawerAdapter
import kotlinx.coroutines.runBlocking
class ChatActivity : AppCompatActivity() {
    private lateinit var binding: ActivityChatBinding
    private lateinit var chatAdapter: ChatAdapter
    private val viewModel: ChatViewModel by viewModels()
    private val message = mutableListOf<Message>()
    private lateinit var userPreference : UserPreference
    private lateinit var drawerLayout : DrawerLayout
    private lateinit var layoutDrawerBinding: LayoutDrawerBinding
    private val drawerViewModel: DrawerViewModel by viewModels()
    private lateinit var drawerAdapter: DrawerAdapter

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityChatBinding.inflate(layoutInflater)
        setContentView(binding.root)

        userPreference = UserPreference.getInstance(this.dataStore)
        chatAdapter = ChatAdapter(message)
        drawerAdapter = DrawerAdapter()

        binding.recyclerView.layoutManager = LinearLayoutManager(this)
        binding.recyclerView.adapter = chatAdapter
        val drawerView = binding.navigationView.getHeaderView(0)
        layoutDrawerBinding = LayoutDrawerBinding.bind(drawerView)

        val layoutManager = LinearLayoutManager(this)
        layoutDrawerBinding.drawerRv.layoutManager = layoutManager

        layoutDrawerBinding.drawerRv.adapter = drawerAdapter
        binding.chatSendButton.setOnClickListener{
            sendButton()
            checkMessageAvailibility()
        }

        viewModel.chatMessage.observe(this) { response ->
            response?.let {
                val responseResult = it.data?.response.toString()
                addMessage(responseResult, false)
            }
        }
        setTitle()
        setDrawer()
        closeDrawer()
        observeChatData()
    }

    private fun observeTopicData() {
        drawerViewModel.topicResponse.observe(this, Observer{topicResponse ->
            Log.d("TOPIC RESPONSE","$topicResponse")
            drawerAdapter.submitList(topicResponse)
        })
    }

    private fun fetchTopic() {
        runBlocking {
            val token = userPreference.getToken()
            if (token != null) {
                drawerViewModel.getTopic(token)
            }
        }
    }



    private fun checkMessageAvailibility() {
        if(message.isNotEmpty()){
            binding.cardviewTitle.visibility = View.INVISIBLE
            binding.cardTextview.visibility = View.INVISIBLE
            binding.card1.visibility = View.INVISIBLE
            binding.card2.visibility = View.INVISIBLE
            binding.chatTvInspirasi.visibility = View.INVISIBLE

            binding.recyclerView.visibility = View.VISIBLE
        }else{
            binding.cardviewTitle.visibility = View.VISIBLE
            binding.cardTextview.visibility = View.VISIBLE
            binding.card1.visibility = View.VISIBLE
            binding.card2.visibility = View.VISIBLE
            binding.chatTvInspirasi.visibility = View.VISIBLE

            binding.recyclerView.visibility = View.INVISIBLE
        }
    }

    private fun setDrawer(){
        drawerLayout = binding.drawerLayout
        binding.chatMenu.setOnClickListener{
            drawerLayout.openDrawer(GravityCompat.END)
            fetchTopic()
            observeTopicData()
        }
    }
    private fun closeDrawer(){
        layoutDrawerBinding.drawerMenu.setOnClickListener{
            drawerLayout.closeDrawer(GravityCompat.END)
        }
    }
    private fun addMessage(text : String, isUser: Boolean){
        val chatMessage = Message(text, isUser)
        message.add(chatMessage)
        chatAdapter.notifyItemInserted(message.size - 1)
        binding.recyclerView.scrollToPosition(message.size - 1)
        if (isUser) binding.chatEditPrompt.text.clear()
    }

    private fun sendButton(){
        runBlocking {
            val token = userPreference.getToken()
            val prompt = binding.chatEditPrompt.text.toString()
            if(prompt.isNotEmpty()){
                addMessage(prompt, true)
                if(token != null){
                    viewModel.getChatResponse(token, prompt)
                }
            }
        }

    }
    private fun observeChatData(){
        viewModel.isLoading.observe(this){isLoading->
            if(isLoading) binding.chatProgressBar.visibility = View.VISIBLE else binding.chatProgressBar.visibility = View.INVISIBLE
        }
    }
    private fun setTitle (){
        val title = intent.getStringExtra("title")
        binding.chatTitle.text = title
    }
}

